package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "url-shortener/docs"
	"url-shortener/internal/config"
	"url-shortener/internal/controllers"
)

type Handler struct {
	UrlController *controllers.UrlController
}

func NewHandler(urlController *controllers.UrlController) *Handler {
	return &Handler{urlController}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
