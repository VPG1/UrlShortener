package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"url-shortener/internal/Logger"
	"url-shortener/internal/config"
	"url-shortener/internal/controllers"
	"url-shortener/internal/services"
	"url-shortener/internal/storage/postgresql"
)

func main() {
	cfg := config.MustLoadConfig()
	log := Logger.SetupLogger(cfg.Env)

	log.Info("Start url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug mod enabled")

	pgStorage, err := postgresql.NewStorage(cfg.PostgresServer, log)
	if err != nil {
		log.Error("Failed to connect to postgres", err)
		return
	}
	defer pgStorage.Close()

	urlService := services.NewUrlService(cfg.AliasLen, pgStorage, log)

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	urlController := controllers.NewUrlController(urlService, log)

	router.GET("/", func(c *gin.Context) {
		urls, err := urlController.GetUrls(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, urls)
		}
	})

	router.POST("/", func(c *gin.Context) {
		alias, err := urlController.ShortenURL(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusCreated, gin.H{"alias": c.FullPath() + alias})
		}
	})

	router.GET("/:alias", func(c *gin.Context) {
		url, err := urlController.GetUrlByAlias(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if url == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		} else {
			c.Redirect(http.StatusTemporaryRedirect, url)
		}
	})

	router.DELETE("/", func(c *gin.Context) {
		isUrlDeleted, err := urlController.DeleteAlias(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if !isUrlDeleted {
			c.JSON(http.StatusNoContent, gin.H{"error": "url alias doesn't exist"})
		} else {
			c.JSON(http.StatusAccepted, gin.H{"status": "alias successfully deleted"})
		}
	})

	err = router.Run(fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, cfg.HTTPServer.Port))
	if err != nil {
		log.Error("Failed to start server")
		return
	}

}
