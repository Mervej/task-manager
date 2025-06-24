package model

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskStatus string

const (
	Todo       TaskStatus = "Todo"
	InProgress TaskStatus = "InProgress"
	Done       TaskStatus = "Done"
)
