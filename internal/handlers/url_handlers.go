package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetAllUrls(c *gin.Context) {
	urls, err := h.UrlController.GetUrls(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, urls)
	}
}

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

func (h *Handler) DeleteUrl(c *gin.Context) {
	isUrlDeleted, err := h.UrlController.DeleteAlias(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if !isUrlDeleted {
		c.JSON(http.StatusNoContent, gin.H{"error": "url alias doesn't exist"})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"status": "alias successfully deleted"})
	}
}
