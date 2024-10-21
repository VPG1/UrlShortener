package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "url-shortener/docs"
	"url-shortener/internal/Logger"
	"url-shortener/internal/config"
	"url-shortener/internal/services"
)

type Handler struct {
	AuthService *services.AuthService
	UrlService  *services.UrlService
	Logger      Logger.Logger
}

func NewHandler(authService *services.AuthService, urlService *services.UrlService, logger Logger.Logger) *Handler {
	return &Handler{authService, urlService, logger}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	root := router.Group("/", h.userIdentity)
	{
		root.GET("/", h.GetAllUrls)

		root.POST("/", h.ShortenUrl)

		root.GET("/:alias", h.Redirect)

		root.DELETE("/", h.DeleteUrl)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sign_up", h.SignUp)
		auth.GET("sign_in", h.SignIn)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
