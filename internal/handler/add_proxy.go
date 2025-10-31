package handler

import (
	"net/http"
	"regexp"

	"github.com/d1manpro/nginx-proxy-api/internal/certbot"
	"github.com/d1manpro/nginx-proxy-api/internal/config"
	"github.com/d1manpro/nginx-proxy-api/internal/nginx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddProxy(cfg *config.Config, log *zap.Logger) func(c *gin.Context) {
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

		err := certbot.GetCert(req.Domain, cfg.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "certbot error"})
			log.Error("failed to get certificate", zap.String("domain", req.Domain), zap.Error(err))
			return
		}

		nginxCfgPath := req.Domain + ".conf"

		err = nginx.AddConfig(req.Domain, req.Target, cfg.NginxCfgTemplate, nginxCfgPath)
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
