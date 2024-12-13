package delivery

import (
	"github.com/gin-gonic/gin"
	"task-queue/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/tasks")
	{
		api.POST("", h.createTask)
		api.GET(":id", h.getTask)
	}
	return router
}
