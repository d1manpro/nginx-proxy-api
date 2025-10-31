package handler

import (
	"net/http"

	"github.com/d1manpro/nginx-proxy-api/internal/certbot"
	"github.com/d1manpro/nginx-proxy-api/internal/config"
	"github.com/d1manpro/nginx-proxy-api/internal/nginx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RemoveProxy(cfg *config.Config, log *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req RemoveDomainReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid JSON: " + err.Error()})
			return
		}

		if req.Domain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "domain is empty"})
			return
		}

		if !isDomainValid(req.Domain) {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "domain is invalid"})
			return
		}

		nginxCfgPath := req.Domain + ".conf"

		err := nginx.RemoveConfig(nginxCfgPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to setup nginx config"})
			log.Error("failed to setup nginx config", zap.String("domain", req.Domain), zap.Error(err))
			return
		}

		err = certbot.DeleteCert(req.Domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "certbot error"})
			log.Error("failed to get certificate", zap.String("domain", req.Domain), zap.Error(err))
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"status": "deleted"})
	}
}
