package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader      = "Authorization"
	userCtx         = "userid"
	petitionIdParam = "id"
	subIdParam      = "sign_id"
)

func (h *Handler) authRequired(c *gin.Context) {
	token := c.GetHeader(authHeader)
	if token == "" {
		newErrorResponse(c, http.StatusUnauthorized, "not authorized")
		return
	}

	headersParts := strings.Split(token, " ")
	if len(headersParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "not correct auth token format")
		return
	}

	userId, err := h.services.ParseToken(headersParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func (h *Handler) getUserId(c *gin.Context) (int, error) {
	userId, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "no authenticated user id found")
		return 0, nil
	}

	userIdInt, ok := userId.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "not correct authenticated user id format")
		return 0, nil
	}

	return userIdInt, nil
}
