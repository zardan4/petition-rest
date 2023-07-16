package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	petitions "github.com/zardan4/petition-rest"
)

type getAllSubsResponses struct {
	Data []petitions.Sub `json:"data"`
}

func (h *Handler) getAllSubs(ctx *gin.Context) {
	petitionId, err := strconv.Atoi(ctx.Param(petitionIdParam))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	subs, err := h.services.Subs.GetAllSubs(petitionId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getAllSubsResponses{
		subs,
	})
}

func (h *Handler) createSub(ctx *gin.Context) {
	userId, err := h.getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	petitionId, err := strconv.Atoi(ctx.Param(petitionIdParam))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	createdSubId, err := h.services.Subs.CreateSub(petitionId, userId)
	if err != nil {
		// check if user already voted
		pqErr, ok := err.(*pq.Error)
		if !ok {
			newErrorResponse(ctx, http.StatusInternalServerError,
				"error while setting custom error type")
			return
		}
		if pqErr.Code.Name() == "unique_violation" {
			newErrorResponse(ctx, http.StatusForbidden, "user already voted")
			return
		}

		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": createdSubId,
	})
}

func (h *Handler) deleteSub(ctx *gin.Context) {
	userId, err := h.getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	petitionId, err := strconv.Atoi(ctx.Param(petitionIdParam))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	subId, err := strconv.Atoi(ctx.Param(subIdParam))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.DeleteSub(subId, petitionId, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) checkSignature(ctx *gin.Context) {
	userId, err := h.getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	petitionId, err := strconv.Atoi(ctx.Param(petitionIdParam))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	signed, err := h.services.CheckSignature(petitionId, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"signed": signed,
	})
}
