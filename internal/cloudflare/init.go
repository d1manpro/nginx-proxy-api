package cloudflare

import (
	"net"
	"net/http"
	"time"

	"github.com/d1manpro/nginx-proxy-api/internal/config"
	"go.uber.org/zap"
)

type CfAPI struct {
	client *http.Client
	cfg    *config.Config
	log    *zap.Logger
}

func InitCfAPI(cfg *config.Config, log *zap.Logger) *CfAPI {
	return &CfAPI{
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &RoundTripper{
				logger: log,
				cfg:    cfg,
				rt: &http.Transport{
					DialContext: (&net.Dialer{
						Timeout:   5 * time.Second,
						KeepAlive: 30 * time.Second,
					}).DialContext,
					TLSHandshakeTimeout: 5 * time.Second,
					MaxIdleConns:        100,
					IdleConnTimeout:     90 * time.Second,
				},
			},
		},
		cfg: cfg,
		log: log,
	}
}

type RoundTripper struct {
	logger *zap.Logger
	cfg    *config.Config
	rt     http.RoundTripper
}

func (l *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+l.cfg.Cloudflare.Token)
	req.Header.Add("Content-Type", "application/json")

	start := time.Now()
	resp, err := l.rt.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		l.logger.Error("HTTP request failed",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return nil, err
	}

	l.logger.Info("HTTP request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Int("status", resp.StatusCode),
		zap.Duration("duration", duration),
	)
	return resp, nil
}
