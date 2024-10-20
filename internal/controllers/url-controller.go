package controllers

import (
	"github.com/gin-gonic/gin"
	"url-shortener/internal/Logger"
	"url-shortener/internal/services"
)

type URLDto struct {
	Url string `json:"url" binding:"required,url"`
}

type UrlController struct {
	Service *services.UrlService
	Logger  Logger.Logger
}

func NewUrlController(service *services.UrlService, logger Logger.Logger) *UrlController {
	return &UrlController{Service: service, Logger: logger}
}

func (urlC *UrlController) ShortenURL(c *gin.Context) (string, error) {
	var urlDto URLDto
	if err := c.ShouldBindJSON(&urlDto); err != nil {
		urlC.Logger.Error("Incorrect body format", "err", err.Error())
		return "", err
	}

	url, err := urlC.Service.CreateNewAlias(urlDto.Url)
	if err != nil {
		urlC.Logger.Error("Incorrect url", "err", err.Error())
		return "", err
	}

	urlC.Logger.Info("Created new url", "url", url.Url, "alias", url.Alias)

	return url.Alias, nil
}
