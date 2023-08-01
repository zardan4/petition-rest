package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	petitions "github.com/zardan4/petition-rest"
)

// @Summary SignUp
// @Tags auth
// @Description Create account
// @ID signup
// @Accept  json
// @Produce  json
// @Param input body petitions.User true "Account info"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/signup [post]
func (h *Handler) signUp(c *gin.Context) {
	var input petitions.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, idResponse{
		Id: id,
	})
}

type singInInput struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInResponse struct {
	Token string `json:"token"`
}

// @Summary SignIn
// @Tags auth
// @Description Enter account
// @ID signin
// @Accept  json
// @Produce  json
// @Param input body singInInput true "Account info"
// @Success 200 {object} signInResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/signin [post]
func (h *Handler) signIn(c *gin.Context) {
	var input singInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Name, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}
