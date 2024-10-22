package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type URLDto struct {
	Url string `json:"url" binding:"required,url"`
}

type AliasDto struct {
	Alias string `json:"alias" binding:"required"`
}

func GetUserId(c *gin.Context) uint64 {
	userIdStr, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewResponseError("User id not found"))
		return 0
	}

	// try convert userId to uint64
	userId, ok := userIdStr.(uint64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewResponseError("Incorrect user id"))
	}

	return userId
}

// @Summary GetAllUserUrls
// @Security ApiKeyAuth
// @Tags urls
// @ID get-all-user-urls
// @Description get all urls
// @Accept json
// @Produce json
// @Success 200 {object} []entities.URl
// @Failure 500 {object} ResponseError
// @Router /api [get]
func (h *Handler) GetAllUserUrls(c *gin.Context) {
	userId := GetUserId(c)

	urls, err := h.UrlService.GetUserUrls(userId)

	if err != nil {
		h.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, NewResponseError(err.Error()))
	} else {
		c.JSON(http.StatusOK, urls)
	}
}

// @Summary ShortenUrl
// @Security ApiKeyAuth
// @Tags urls
// @Description give alias for url
// @ID shorten-url
// @Accept json
// @Produce json
// @Param input body URLDto true "link"
// @Success 200 {object} URLDto
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /api [post]
func (h *Handler) ShortenUrl(c *gin.Context) {
	userId := GetUserId(c)

	var urlDto URLDto
	if err := c.ShouldBindJSON(&urlDto); err != nil { // возможно не корректный url
		c.JSON(http.StatusBadRequest, NewResponseError("Incorrect body format"))
		return
	}

	url, err := h.UrlService.CreateNewAlias(urlDto.Url, userId)
	if err != nil {
		h.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, NewResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"alias": c.Request.Host + "/" + url.Alias})
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
// @Security ApiKeyAuth
// @Tags urls
// @Description delete alias
// @ID delete-url
// @Accept json
// @Produce json
// @Param input body AliasDto false "alias"
// @Success 200 {object} SuccessResponse
// @Success 204 {object} SuccessResponse
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /api [delete]
func (h *Handler) DeleteUrl(c *gin.Context) {
	userId := GetUserId(c)

	var alias AliasDto
	if err := c.ShouldBindJSON(&alias); err != nil {
		h.Logger.Error("Incorrect body format", "err", err.Error())
		c.JSON(http.StatusBadRequest, NewResponseError("Incorrect body format"))
		return
	}

	isUrlDeleted, err := h.UrlService.DeleteUrlByAlias(alias.Alias, userId)
	if err != nil {
		h.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, NewResponseError(err.Error()))
	} else if !isUrlDeleted {
		c.JSON(http.StatusNoContent, NewSuccessResponse("url alias doesn't exist"))
	} else {
		c.JSON(http.StatusOK, NewSuccessResponse("alias successfully deleted"))
	}
}
