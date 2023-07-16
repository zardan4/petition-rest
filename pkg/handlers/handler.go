package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zardan4/petition-rest/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
	}

	api := router.Group("/api", h.authRequired)
	{
		petitions := api.Group("/petitions")
		{
			petitions.GET("/", h.getAllPetitions)
			petitions.POST("/", h.createPetition)

			petitions.GET("/:id", h.getPetition)
			petitions.PUT("/:id", h.updatePetition)
			petitions.DELETE("/:id", h.deletePetition)

			petitions.GET("/:id/signed", h.checkSignature) // check if user already signed petition

			signs := petitions.Group("/:id/subs")
			{
				signs.GET("/", h.getAllSubs)
				signs.POST("/", h.createSub)

				signs.DELETE("/:sign_id", h.deleteSub)
			}
		}
	}

	return router
}
