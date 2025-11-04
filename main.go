package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/d1manpro/nginx-proxy-api/internal/cloudflare"
	"github.com/d1manpro/nginx-proxy-api/internal/config"
	"github.com/d1manpro/nginx-proxy-api/internal/router"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log := setupLogger()
	defer log.Sync()
	log.Info("Starting...\n")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	cfAPI := cloudflare.InitCfAPI(cfg, log)

	server := router.NewServer(cfg, log, cfAPI)
	server.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	server.Stop()
	log.Info("HTTP-server stopped")

	log.Info("Script stopped")
}

func setupLogger() *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeTime:  zapcore.TimeEncoderOfLayout("2006.01.02 15:04:05.000"),
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	return zap.New(zapcore.NewTee(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)))
}
