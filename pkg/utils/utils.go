package utils

import (
	"task-manager/internal/model"
)

func IsValidTaskStatus(status string) bool {
	switch model.TaskStatus(status) {
	case model.Todo, model.InProgress, model.Done:
		return true
	default:
		return false
	}
}
