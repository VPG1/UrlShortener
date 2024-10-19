package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"url-shortener/internal/Logger"
	"url-shortener/internal/config"
	"url-shortener/internal/services"
	"url-shortener/internal/storage/postgresql"
)

func main() {
	cfg := config.MustLoadConfig()
	log := Logger.SetupLogger(cfg.Env)

	log.Info("Start url-shortener", slog.String("env", cfg.Env))
	log.Debug("enabled")

	pgStorage, err := postgresql.NewStorage(cfg.PostgresServer, log)
	if err != nil {
		log.Error("Failed to connect to postgres", err)
		return
	}
	defer pgStorage.Close()

	urlService := services.NewUrlService(8, pgStorage, log)
	url, _ := urlService.GetUrlByAlias("TLFpKHNp")
	fmt.Println(url)
	fmt.Println(pgStorage.SelectAll())
}
