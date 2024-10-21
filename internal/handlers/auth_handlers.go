package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Username string `json:"username" binding:"required,min=8,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, NewResponseError("Incorrect body format"))
		return
	}

	user, err := h.AuthService.CreateUser(input.Name, input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewResponseError("Failed to create user"))
		return
	}

	c.JSON(http.StatusCreated, NewSuccessResponse(user.UserName+" "+user.Name))
}

type SignInInput struct {
	Username string `json:"username" binding:"required,min=8,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, NewResponseError("Incorrect body format"))
		return
	}

	token, err := h.AuthService.GenerateToken(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewResponseError("Failed to create token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
