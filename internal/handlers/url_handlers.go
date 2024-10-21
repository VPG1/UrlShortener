package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type URLDto struct {
	Url string `json:"url" binding:"required,url"`
}

type AliasDto struct {
	Alias string `json:"alias" binding:"required"`
}

// @Summary GetAllUrls
// @Tags urls
// @Description get all urls
// @Accept json
// @Produce json
// @Success 200 {object} []string "qwer"
// @Failure 500 {object} ResponseError
// @Router / [get]
func (h *Handler) GetAllUrls(c *gin.Context) {
	urls, err := h.UrlService.GetUrls()

	if err != nil {
		h.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, NewResponseError(err.Error()))
	} else {
		c.JSON(http.StatusOK, urls)
	}
}

// @Summary ShortenUrl
// @Tags urls
// @Description give alias for url
// @Accept json
// @Produce json
// @Param input body URLDto true "link"
// @Success 200 {object} URLDto
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router / [post]
func (h *Handler) ShortenUrl(c *gin.Context) {
	var urlDto URLDto
	if err := c.ShouldBindJSON(&urlDto); err != nil { // возможно не корректный url
		c.JSON(http.StatusBadRequest, NewResponseError("Incorrect body format"))
		return
	}

	url, err := h.UrlService.CreateNewAlias(urlDto.Url)
	if err != nil {
		h.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, NewResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"alias": c.FullPath() + url.Alias})
}

func (h *Handler) Redirect(c *gin.Context) {
	alias := c.Param("alias")

	url, err := h.UrlService.GetUrlByAlias(alias)
	if err != nil {
		h.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, NewResponseError(err.Error()))
	} else if url == nil {
		c.JSON(http.StatusNotFound, NewResponseError("Url not found"))
	} else {
		c.Redirect(http.StatusTemporaryRedirect, url.Url)
	}
}

// @Summary DeleteUrl
// @Tags urls
// @Description delete alias
// @Accept json
// @Produce json
// @Param input body AliasDto false "alias"
// @Success 200 {object} SuccessResponse
// @Success 204 {object} SuccessResponse
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router / [delete]
func (h *Handler) DeleteUrl(c *gin.Context) {
	var alias AliasDto
	if err := c.ShouldBindJSON(&alias); err != nil {
		h.Logger.Error("Incorrect body format", "err", err.Error())
		c.JSON(http.StatusBadRequest, NewResponseError("Incorrect body format"))
		return
	}

	isUrlDeleted, err := h.UrlService.DeleteUrlByAlias(alias.Alias)
	if err != nil {
		h.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, NewResponseError(err.Error()))
	} else if !isUrlDeleted {
		c.JSON(http.StatusNoContent, NewSuccessResponse("url alias doesn't exist"))
	} else {
		c.JSON(http.StatusOK, NewSuccessResponse("alias successfully deleted"))
	}
}
