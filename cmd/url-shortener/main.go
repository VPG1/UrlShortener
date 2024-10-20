package main

import (
	"context"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

	log.Debug("Postgres connection established")

	// creating service
	urlService := services.NewUrlService(cfg.AliasLen, pgStorage, log)

	// initialize routes and start server
	urlController := controllers.NewUrlController(urlService, log)
	handler := handlers.NewHandler(urlController)

	router := handler.InitRoutes(cfg, urlService, log)

	// setting up http server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router.Handler(),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to listen and server", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown:", err)
		panic(err)
	}

	longShutdown := make(chan struct{}, 1)

	go func() {
		pgStorage.Close()
		longShutdown <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		log.Info("Timeout of long shutdown")
	case <-longShutdown:
		log.Info("finished")
	}
}
