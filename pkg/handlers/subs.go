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

// @Summary Gets all signatories
// @Security ApiKeyAuth
// @Tags signatories
// @Description Get all signatories by petition
// @ID get-signatories
// @Accept  json
// @Produce  json
// @Param id path int true "Petition id"
// @Success 200 {object} getAllSubsResponses
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions/{id}/subs [get]
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

// @Summary Create signatorie
// @Security ApiKeyAuth
// @Tags signatories
// @Description Creates new signatorie by petition
// @ID create-signatorie
// @Accept  json
// @Produce  json
// @Param id path int true "Petition id"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions/{id}/subs [post]
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

		if pqErr.Code.Name() != "unique_violation" {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		newErrorResponse(ctx, http.StatusForbidden, "user already voted")
		return
	}

	ctx.JSON(http.StatusCreated, idResponse{
		Id: createdSubId,
	})
}

// @Summary Deletes signatorie
// @Security ApiKeyAuth
// @Tags signatories
// @Description Deletes signatorie by petition
// @ID delete-signatorie
// @Accept  json
// @Produce  json
// @Param id path int true "Petition id"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions/{id}/subs [delete]
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

	err = h.services.DeleteSub(petitionId, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

type checkSignatureResponse struct {
	Signed bool `json:"signed"`
}

// @Summary Checks signatorie
// @Security ApiKeyAuth
// @Tags signatories
// @Description Checks signatorie by petition
// @ID check-signatorie
// @Accept  json
// @Produce  json
// @Param id path int true "Petition id"
// @Success 200 {object} checkSignatureResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions/{id}/signed/ [get]
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

	ctx.JSON(http.StatusOK, checkSignatureResponse{
		Signed: signed,
	})
}
