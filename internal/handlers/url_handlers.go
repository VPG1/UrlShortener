package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

func NewResponseError(error string) ResponseError {
	return ResponseError{Error: error}
}

// @Summary GetAllUrls
// @Tags urls
// @Description get all urls
// @Accept json
// @Produce json
// @Success 200 {object} []string{}
// @Failure 400 {object} ResponseError
// @Router /urls [get]

func (h *Handler) GetAllUrls(c *gin.Context) {
	urls, err := h.UrlController.GetUrls(c)
	if err != nil {
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
// @Param input body controllers.URLDto false "link"
// @Success 200 {object} controllers.URLDto
// @Failure 400 {object} ResponseError
// @Router /urls [post]

func (h *Handler) ShortenUrl(c *gin.Context) {
	alias, err := h.UrlController.ShortenURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, gin.H{"alias": c.FullPath() + alias})
	}
}

func (h *Handler) Redirect(c *gin.Context) {
	url, err := h.UrlController.GetUrlByAlias(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if url == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
	} else {
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// @Summary DeleteUrl
// @Tags urls
// @Description delete alias
// @Accept json
// @Produce json
// @Param input body controllers.AliasDto false "alias"
// @Success 200 {object} controllers.URLDto
// @Success 204 {object} ResponseError
// @Failure 400 {object} ResponseError
// @Router /urls [delete]

func (h *Handler) DeleteUrl(c *gin.Context) {
	isUrlDeleted, err := h.UrlController.DeleteAlias(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if !isUrlDeleted {
		c.JSON(http.StatusNoContent, gin.H{"error": "url alias doesn't exist"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "alias successfully deleted"})
	}
}
