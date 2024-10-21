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

type SignUpOutput struct {
	Status string `json:"status"`
	Id     uint64 `json:"id"`
}

func NewSuccessRegistrationResponse(status string, id uint64) *SignUpOutput {
	return &SignUpOutput{status, id}
}

// @Summary SignUp
// @Tags auth
// @ID sign-up
// @Description Registration
// @Accept json
// @Produce json
// @Param input body SignUpInput true "SignUp input"
// @Success 201 {object} SignUpOutput
// @Failure 401 {object} ResponseError
// @Router /auth/sign_up [post]
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

	c.JSON(http.StatusCreated, NewSuccessRegistrationResponse("User successfully created", user.Id))
}

type SignInInput struct {
	Username string `json:"username" binding:"required,min=8,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type SignInOutput struct {
	Token string `json:"token"`
}

func NewSuccessTokenGenerationResponse(token string) *SignInOutput {
	return &SignInOutput{token}
}

// @Summary SignIn
// @Tags auth
// @ID sign-in
// @Description Generate token
// @Accept json
// @Produce json
// @Param input body SignInInput true "SignIn input"
// @Success 200 {object} SignInOutput
// @Failure 401 {object} ResponseError
// @Router /auth/sign_in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, NewResponseError("Incorrect body format"))
		return
	}

	token, err := h.AuthService.GenerateToken(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewResponseError("Failed to create token: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessTokenGenerationResponse(token))
}
