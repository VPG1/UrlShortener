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

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
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

	router.GET("/:alias", h.Redirect)

	root := router.Group("/api", h.userIdentity)
	{
		root.GET("/", h.GetAllUserUrls)

		root.POST("/", h.ShortenUrl)

		root.DELETE("/", h.DeleteUrl)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sign_up", h.SignUp)
		auth.POST("sign_in", h.SignIn)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
