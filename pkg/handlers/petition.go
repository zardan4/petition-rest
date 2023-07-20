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

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

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
