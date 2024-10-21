package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const authorizationHeader = "Authorization"

func (h *Handler) userIdentity(c *gin.Context) {
	authHeader := c.GetHeader(authorizationHeader)
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponseError("empty authorization header"))
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponseError("invalid authorization header"))
		return
	}

	userId, err := h.AuthService.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponseError("invalid token"))
		return
	}

	c.Set("userId", userId)
}
