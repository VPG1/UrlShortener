package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"url-shortener/internal/Logger"
	"url-shortener/internal/config"
	"url-shortener/internal/controllers"
	"url-shortener/internal/handlers"
	"url-shortener/internal/services"
	"url-shortener/internal/storage/postgresql"
)

// @title URL Shortener App
// @version 1.0
// @description API Service for URL shorten

// @host localhost:8080
// @BasePath /

func main() {
	// load config
	cfg := config.MustLoadConfig()

	// setup logger
	log := Logger.SetupLogger(cfg.Env)

	log.Info("Start url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug mod enabled")

	// create storage
	pgStorage, err := postgresql.NewStorage(cfg.PostgresServer, log)
	if err != nil {
		log.Error("Failed to connect to postgres", err)
		return
	}
	defer pgStorage.Close()

	log.Debug("Postgres connection established")

	// creating service
	urlService := services.NewUrlService(cfg.AliasLen, pgStorage, log)

	// initialize routes and start server
	urlController := controllers.NewUrlController(urlService, log)
	handler := handlers.NewHandler(urlController)

	router := handler.InitRoutes(cfg, urlService, log)
	err = router.Run(fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, cfg.HTTPServer.Port))
	if err != nil {
		log.Error("Failed to start server")
		return
	}
}
