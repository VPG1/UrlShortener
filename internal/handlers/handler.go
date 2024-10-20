package handlers

import (
	"github.com/gin-gonic/gin"
	"url-shortener/internal/Logger"
	"url-shortener/internal/config"
	"url-shortener/internal/controllers"
	"url-shortener/internal/services"
)

type Handler struct {
	UrlController *controllers.UrlController
}

func NewHandler(urlController *controllers.UrlController) *Handler {
	return &Handler{urlController}
}

func (h *Handler) InitRoutes(cfg *config.Config, urlService *services.UrlService, log Logger.Logger) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	root := router.Group("/")
	{
		root.GET("/", h.GetAllUrls)

		root.POST("/", h.ShortenUrl)

		root.GET("/:alias", h.Redirect)

		root.DELETE("/", h.DeleteUrl)
	}

	return router
}
