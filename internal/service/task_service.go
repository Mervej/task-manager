package service

import (
	"errors"
	"strconv"
	"task-manager/internal/model"
	"task-manager/internal/repository"
)

type TaskService struct {
	repo repository.ITaskRepository
}

func NewTaskService(repo repository.ITaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *model.Task, userIDStr string) (*model.Task, error) {
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}

	if task.Title == "" {
		return nil, errors.New("task title is required")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	task.UserID = userID

	return s.repo.Create(task)
}

func (s *TaskService) GetTasks(userIDStr string, status ...string) ([]model.Task, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// check if any status filter is provided
	var statusFilter string
	if len(status) > 0 {
		statusFilter = status[0]
	}
	return s.repo.GetAll(userID, statusFilter)
}

func (s *TaskService) GetTaskByID(id int, userIDStr string) (*model.Task, error) {
	if id <= 0 {
		return nil, errors.New("invalid task ID")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByID(id, userID)
}

func (s *TaskService) UpdateTask(id int, task *model.Task, userIDStr string) (*model.Task, error) {
	if id <= 0 {
		return nil, errors.New("invalid task ID")
	}
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// check if the task exists for the user
	existingTask, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	// update the given fields
	existingTask.Title = task.Title
	existingTask.Description = task.Description
	existingTask.DueDate = task.DueDate

	// only update status if provided and valid
	if task.Status != "" {
		if task.Status != string(model.Todo) &&
			task.Status != string(model.InProgress) &&
			task.Status != string(model.Done) {
			return nil, errors.New("invalid task status")
		}
		existingTask.Status = task.Status
	}

	return s.repo.Update(existingTask)
}

func (s *TaskService) DeleteTask(id int, userIDStr string) (string, error) {
	if id <= 0 {
		return "", errors.New("invalid task ID")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return "", errors.New("invalid user ID")
	}

	return s.repo.Delete(id, userID)
}
