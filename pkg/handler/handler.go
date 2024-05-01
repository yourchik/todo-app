package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourchik/todo-app/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.singUp)
		auth.POST("sing-in", h.singIn)
	}

	api := router.Group("/api", h.userIndentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
				items.GET("/:id", h.getListById)
				items.PUT("/:id", h.updateItem)
				items.DELETE("/:id", h.deleteItem)
			}
		}
	}
	return router
}
