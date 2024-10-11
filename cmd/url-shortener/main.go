package main

import (
	"log/slog"
	"url-shortener/internal/Logger"
	"url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()
	log := Logger.SetupLogger(cfg.Env)

	log.Info("Start url-shortener", slog.String("env", cfg.Env))
	log.Debug("enabled")
}
