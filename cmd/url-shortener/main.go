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

	// TODO: move alias len to config
	urlService := services.NewUrlService(8, pgStorage, log)

	router := gin.Default()

	urlController := controllers.NewUrlController(urlService, log)

	router.POST("/", func(c *gin.Context) {
		alias, err := urlController.ShortenURL(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusCreated, gin.H{"alias": c.FullPath() + alias})
		}
	})

	router.GET("/:alias", func(c *gin.Context) {
		alias, err := urlController.GetUrlByAlias(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.Redirect(http.StatusTemporaryRedirect, alias)
		}
	})

	err = router.Run(fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, cfg.HTTPServer.Port))
	if err != nil {
		log.Error("Failed to start server")
		return
	}

}
