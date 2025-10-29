package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	server := router.NewServer(cfg, log)
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

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	var cores []zapcore.Core
	cores = append(cores, consoleCore)

	logFile, err := os.OpenFile("pter-api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("cannot open log file: %v", err))
	}
	fileCore := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.InfoLevel)
	cores = append(cores, fileCore)

	return zap.New(zapcore.NewTee(cores...))
}
