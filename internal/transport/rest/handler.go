package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zardan4/petition-rest/docs"
	"github.com/zardan4/petition-rest/internal/service"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
		auth.POST("/refresh", h.refreshTokens)
	}

	api := router.Group("/api", h.authRequired, h.userExists)
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

				signs.DELETE("/", h.deleteSub)
			}
		}
	}

	return router
}
