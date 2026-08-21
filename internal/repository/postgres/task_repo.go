package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"task-manager/internal/domain"
	castomErr "task-manager/pkg/errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, userID uuid.UUID, req *domain.CreateTaskRequest) (*domain.Task, error) {
	var task domain.Task
	query := `
		INSERT INTO task (user_id, title, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, description, status, create_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, userID, req.Title, req.Description, req.Status).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("task_repo.Create: %w", err)
	}

	return &task, nil
}

func (r *TaskRepo) GetByID(ctx context.Context, userID, taskID uuid.UUID) (*domain.Task, error) {
	var task *domain.Task
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM task
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.GetContext(ctx, &task, query, taskID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, castomErr.ErrTaskNotFound
		}
		return nil, fmt.Errorf("task_repo.GetByID: %w", err)
	}
	return task, nil
}

func (r *TaskRepo) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Task, int, error) {
	var (
		tasks []domain.Task
		total int
	)
	query := `
		SELECT COUNT(*) 
		FROM tasks
		WHERE user_id = $1
	`
	err := r.db.GetContext(ctx, &total, query, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("task_repo.GetByUserID: %w", err)
	}
	query = `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &tasks, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("task_repo.GetByUserID: %w", err)
	}
	return tasks, total, nil
}

func (r *TaskRepo) Update(ctx context.Context, userID, taskID uuid.UUID, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	var task domain.Task
	now := time.Now()
	query := `
		UPDATE tasks
		SET title = $3, description = $4, status = $5, updated_at = $6
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, title, description, status, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, taskID, userID, req.Title, req.Description, req.Status, now).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, castomErr.ErrTaskNotFound
		}
		return nil, fmt.Errorf("task_repo.Update: %w", err)
	}
	return &task, nil
}

func (r *TaskRepo) Delete(ctx context.Context, userID, taskID uuid.UUID) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.ExecContext(ctx, query, taskID, userID)
	if err != nil {
		return fmt.Errorf("task_repo.Delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("task_repo.Delete: rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return castomErr.ErrTaskNotFound
	}
	return nil
}
