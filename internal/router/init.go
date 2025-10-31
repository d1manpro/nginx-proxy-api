package router

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/d1manpro/nginx-proxy-api/internal/config"
	"github.com/d1manpro/nginx-proxy-api/internal/handler"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server
	cfg    *config.Config
	log    *zap.Logger
}

func NewServer(cfg *config.Config, log *zap.Logger) *Server {
	if !cfg.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.SetTrustedProxies([]string{cfg.Server.Host})

	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))

	r.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !slices.Contains(cfg.Access.AllowedIPs, clientIP) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})

	/*r.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != cfg.Access.Token {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	})*/

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.Origins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	return &Server{
		router: r,
		cfg:    cfg,
		log:    log,
	}
}

func (s *Server) Start() {
	s.router.GET("/test")
	s.router.POST("/add-proxy", handler.AddProxy(s.cfg, s.log))
	s.router.POST("/remove-proxy", handler.RemoveProxy(s.cfg, s.log))

	port := ":" + s.cfg.Server.Port

	s.srv = &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	s.log.Info("Starting server", zap.String("port", port))

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal("server failed", zap.Error(err))
		}
	}()
}

func (s *Server) Stop() {
	if s.srv == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Error("server shutdown error", zap.Error(err))
	}
}
