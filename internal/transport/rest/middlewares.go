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
	// no token was find
	if token == "" {
		newErrorResponse(c, http.StatusUnauthorized, "not correct auth token header")
		return
	}

	headersParts := strings.Split(token, " ")
	// if not enough parts of token
	if len(headersParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "not correct auth token value")
		return
	}

	// if not correct "Bearer" part
	if headersParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "not correct auth token value")
		return
	}

	// if no token
	if len(strings.Trim(headersParts[1], " ")) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := h.services.ParseToken(headersParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

// func (h *Handler) checkUsersPetition(ctx *gin.Context) {
// 	userId, err := h.getUserId(ctx)
// 	if err != nil {
// 		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	petitionId, err := strconv.Atoi(ctx.Param(petitionIdParam))

// 	petition, err := h.services.GetPetition(petitionId)
// 	if err != nil {
// 		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	if petition.
// }

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
