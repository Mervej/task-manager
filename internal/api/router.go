package api

import (
	"task-manager/internal/api/handler"
	"task-manager/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	router := gin.New()

	router.Use(middleware.LoggingMiddleware())

	// initialise the task routes with auth middleware
	tasks := router.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("", taskHandler.GetTasks)
		tasks.GET("/:id", taskHandler.GetTaskByID)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}

	return router
}
