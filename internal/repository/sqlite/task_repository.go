package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"task-manager/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) Create(task *model.Task) (*model.Task, error) {
	query := `
		INSERT INTO tasks (title, description, user_id, due_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now

	result, err := r.db.ExecContext(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.UserID,
		task.DueDate,
		string(model.Todo),
		now,
		now,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	task.ID = int(id)

	return task, nil
}

func (r *TaskRepository) GetAll(userID int, status model.TaskStatus) ([]model.Task, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, title, description, user_id, due_date, status, created_at, updated_at
			FROM tasks
			WHERE user_id = ? AND status = ?
			ORDER BY created_at DESC
		`
		args = append(args, userID, status)
	} else {
		query = `
			SELECT id, title, description, user_id, due_date, status, created_at, updated_at
			FROM tasks
			WHERE user_id = ?
			ORDER BY created_at DESC
		`
		args = append(args, userID)
	}

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.UserID,
			&dueDate,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning task: %w", err)
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(id int, userID int) (*model.Task, error) {
	query := `
		SELECT id, title, description, user_id, due_date, status, created_at, updated_at
		FROM tasks
		WHERE id = ? AND user_id = ?
	`

	row := r.db.QueryRowContext(context.Background(), query, id, userID)

	var task model.Task
	var dueDate sql.NullTime

	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.UserID,
		&dueDate,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found or not accessible: %d", id)
		}
		return nil, fmt.Errorf("error getting task: %w", err)
	}

	if dueDate.Valid {
		task.DueDate = dueDate.Time
	}

	return &task, nil
}

func (r *TaskRepository) Update(task *model.Task) (*model.Task, error) {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, due_date = ?, status = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	now := time.Now().UTC()
	task.UpdatedAt = now

	result, err := r.db.ExecContext(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Status,
		now,
		task.ID,
		task.UserID,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("task not found or not accessible: %d", task.ID)
	}

	return task, nil
}

func (r *TaskRepository) Delete(id int, userID int) (string, error) {
	query := `DELETE FROM tasks WHERE id = ? AND user_id = ?`

	result, err := r.db.ExecContext(context.Background(), query, id, userID)
	if err != nil {
		return "", fmt.Errorf("error deleting task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("task not found or not accessible: %d", id)
	}

	return "Successfully deleted task", nil
}
