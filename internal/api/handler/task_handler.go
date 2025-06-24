package handler

import (
	"net/http"
	"strconv"
	"task-manager/internal/model"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user id in request header"})
		return
	}

	// bind the request body to the task model
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTask, err := h.taskService.CreateTask(&task, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user id in request header"})
		return
	}

	status := c.Query("status")

	tasks, err := h.taskService.GetTasks(userID.(string), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user id in request header"})
		return
	}

	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID.(string))
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or you don't have permission to view it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user id in request header"})
		return
	}

	id := c.Param("id")
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	updatedTask, err := h.taskService.UpdateTask(taskID, &task, userID.(string))
	if err != nil {
		if err.Error() == "task not found or not accessible: "+id {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or you don't have permission to modify it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user id in request header"})
		return
	}

	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	res, err := h.taskService.DeleteTask(taskID, userID.(string))
	if err != nil {
		if err.Error() == "task not found or not accessible: "+id {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or you don't have permission to delete it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
