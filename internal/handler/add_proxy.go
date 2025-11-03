package handler

import (
	"errors"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/d1manpro/nginx-proxy-api/internal/certbot"
	"github.com/d1manpro/nginx-proxy-api/internal/cloudflare"
	"github.com/d1manpro/nginx-proxy-api/internal/config"
	"github.com/d1manpro/nginx-proxy-api/internal/nginx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddProxy(cfg *config.Config, log *zap.Logger, cf *cloudflare.CfAPI) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req AddDomainReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid JSON: " + err.Error()})
			return
		}

		if req.Domain == "" || req.Target == "" {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "domain or target is empty"})
			return
		}

		if !isDomainValid(req.Domain) {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "domain is invalid"})
			return
		}
		if !isTargetValid(req.Target) {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "target is invalid"})
			return
		}

		var zoneID, subdomain, certDomain string
		for k, v := range cfg.Cloudflare.Domains {
			if strings.HasSuffix(req.Domain, k) {
				zoneID = v
				subdomain = strings.TrimSuffix(req.Domain, "."+k)
				certDomain = k
			}
		}

		if zoneID != "" {
			iCE, err := certbot.IsCertExists(certDomain)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "certbot error"})
				log.Error("failed to get certificates list", zap.String("domain", req.Domain), zap.Error(err))
				return
			}
			if !iCE {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "certificate error"})
				log.Warn("certificate not found", zap.String("domain", req.Domain))
				return
			}

			domains, err := cf.GetAllSubdomains(req.Domain, zoneID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "cloudflare error"})
				log.Error("failed to get subdomain list", zap.String("domain", req.Domain), zap.Error(err))
				return
			}
			if slices.Contains(domains, req.Domain) {
				c.JSON(http.StatusConflict, gin.H{"error": "subdomain already taken"})
				return
			}

			err = cf.CreateDNSRecord(zoneID, map[string]any{
				"type":    "A",
				"name":    subdomain,
				"content": cfg.Cloudflare.NodeIP,
				"ttl":     1,
				"proxied": true,
			})
			if err != nil {
				if errors.Is(err, cloudflare.ErrRecordExists) {
					c.JSON(http.StatusConflict, gin.H{"error": "subdomain already taken"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "cloudflare error"})
					log.Error("failed to create dns record", zap.String("domain", req.Domain), zap.Error(err))
				}
				return
			}
		} else {
			certDomain = req.Domain
			err := certbot.GetCert(req.Domain, cfg.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "certbot error"})
				log.Error("failed to get certificate", zap.String("domain", req.Domain), zap.Error(err))
				return
			}
		}

		err := nginx.AddConfig(req.Domain, certDomain, req.Target, cfg.NginxCfgTemplate, req.Domain+".conf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to setup nginx config"})
			log.Error("failed to setup nginx config", zap.String("domain", req.Domain), zap.Error(err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	}
}

func isDomainValid(domain string) bool {
	re := regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	return re.MatchString(domain)
}

func isTargetValid(domain string) bool {
	re := regexp.MustCompile(
		`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?` +
			`(?:\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?){1,2}` +
			`:(?:` +
			`6553[0-5]|655[0-2]\d|65[0-4]\d{2}|6[0-4]\d{3}|[1-5]\d{4}|[1-9]\d{0,3}` +
			`)$`)
	return re.MatchString(domain)
}
