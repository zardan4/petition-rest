package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zardan4/petition-rest/internal/core"
)

const (
	refreshTokenCookie = "refresh_token"
)

// @Summary SignUp
// @Tags auth
// @Description Create account
// @ID signup
// @Accept  json
// @Produce  json
// @Param input body core.User true "Account info"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/signup [post]
func (h *Handler) signUp(c *gin.Context) {
	var input core.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.services.Authorization.CreateUser(c, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, idResponse{
		Id: id,
	})
}

type singInInput struct {
	Name        string `json:"name" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type signInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary SignIn
// @Tags auth
// @Description Enter account
// @ID signin
// @Accept  json
// @Produce  json
// @Param input body singInInput true "Account info and fingerprint"
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

	tokens, err := h.services.Authorization.GenerateTokens(c, input.Name, input.Password, input.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(
		refreshTokenCookie,
		tokens.RefreshToken,
		int(tokens.RefreshTokenTTL.Seconds()),
		"/auth",
		viper.GetString("serverDomain"),
		true,
		true,
	)

	c.JSON(http.StatusOK, signInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type refreshTokensInput struct {
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type refreshTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary Refresh tokens
// @Tags auth
// @Description Refresh pair of tokens
// @ID refresh
// @Accept  json
// @Produce  json
// @Param input body refreshTokensInput true "Fingerprint. and refresh token in cookie"
// @Success 200 {object} refreshTokensResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [post]
func (h *Handler) refreshTokens(c *gin.Context) {
	var input refreshTokensInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no refresh token cookie provided")
		return
	}

	tokens, err := h.services.Authorization.RefreshTokens(refreshToken, input.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(
		refreshTokenCookie,
		tokens.RefreshToken,
		int(tokens.RefreshTokenTTL.Seconds()),
		"/auth",
		viper.GetString("serverDomain"),
		true,
		true,
	)

	c.JSON(http.StatusOK, refreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type logOutInput struct { // nolint:unused
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Log out
// @Tags auth
// @Description Delete refresh session
// @ID logout
// @Accept  json
// @Produce  json
// @Param input body logOutInput true "refresh token. *in cookie"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/logout [post]
func (h *Handler) logOut(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no refresh token cookie provided")
		return
	}

	err = h.services.Authorization.Logout(c, refreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
