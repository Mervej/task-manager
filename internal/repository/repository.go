package repository

import (
	"task-manager/internal/model"
)

type ITaskRepository interface {
	// Create adds a new task to the database
	Create(task *model.Task) (*model.Task, error)

	// GetAll retrieves all tasks for a specific user, optionally filtered by status
	GetAll(userID int, status model.TaskStatus) ([]model.Task, error)

	// GetByID retrieves a task by its ID and verifies it belongs to the user
	GetByID(id int, userID int) (*model.Task, error)

	// Update modifies an existing task
	Update(task *model.Task) (*model.Task, error)

	// Delete removes a task by ID after verifying it belongs to the user
	Delete(id int, userID int) (string, error)
}
