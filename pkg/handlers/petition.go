package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	petitions "github.com/zardan4/petition-rest"
)

type getAllPetitionsResponses struct {
	Data []petitions.Petition `json:"data"`
}

// @Summary Get all petitions
// @Security ApiKeyAuth
// @Tags petitions
// @Description Get all petitions
// @ID get-petitions
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllPetitionsResponses
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions [get]
func (h *Handler) getAllPetitions(ctx *gin.Context) {
	petitions, err := h.services.Petition.GetAllPetitions()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getAllPetitionsResponses{
		petitions,
	})
}

// @Summary Get petition
// @Security ApiKeyAuth
// @Tags petitions
// @Description Get petition
// @ID get-petition
// @Accept  json
// @Produce  json
// @Param id path int true "Petition id"
// @Success 200 {object} petitions.Petition
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions/{id} [get]
func (h *Handler) getPetition(ctx *gin.Context) {
	petitionId, err := strconv.Atoi(ctx.Param(petitionIdParam))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	petition, err := h.services.Petition.GetPetition(petitionId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, petition)
}

type createPetitionInput struct {
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
}

// @Summary Create petition
// @Security ApiKeyAuth
// @Tags petitions
// @Description Creates new petition
// @ID create-petition
// @Accept  json
// @Produce  json
// @Param input body createPetitionInput true "Petition info"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions [post]
func (h *Handler) createPetition(ctx *gin.Context) {
	var input createPetitionInput
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid request body").Error())
		return
	}

	userId, err := h.getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Petition.CreatePetition(input.Title, input.Text, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, idResponse{
		Id: id,
	})
}

// @Summary Delete petition
// @Security ApiKeyAuth
// @Tags petitions
// @Description Delete petition by id
// @ID delete-petition
// @Accept  json
// @Produce  json
// @Param id path int true "Petition id"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions/{id} [delete]
func (h *Handler) deletePetition(ctx *gin.Context) {
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

	err = h.services.Petition.DeletePetition(petitionId, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Update petition
// @Security ApiKeyAuth
// @Tags petitions
// @Description Update petition by id
// @ID update-petition
// @Accept  json
// @Produce  json
// @Param id path int true "Petition id"
// @Param input body petitions.UpdatePetitionInput true "Updated petition content"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/petitions/{id} [put]
func (h *Handler) updatePetition(ctx *gin.Context) {
	var input petitions.UpdatePetitionInput
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid request body").Error())
		return
	}

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

	err = h.services.Petition.UpdatePetition(input, petitionId, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
