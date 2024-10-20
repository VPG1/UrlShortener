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

func (urlC *UrlController) GetUrls(c *gin.Context) ([]string, error) {
	urls, err := urlC.Service.GetUrls()
	if err != nil {
		urlC.Logger.Error(err.Error())
		return nil, err
	}

	return urls, nil
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

func (urlC *UrlController) GetUrlByAlias(c *gin.Context) (string, error) {
	alias := c.Param("alias")
	url, err := urlC.Service.GetUrlByAlias(alias)
	if err != nil {
		urlC.Logger.Error("Incorrect alias", "alias", alias, "err", err.Error())
		return "", err
	}
	urlC.Logger.Info("Get url by alias", "alias", alias, "url", url)

	return url.Url, nil
}
