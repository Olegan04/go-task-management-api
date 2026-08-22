package usecase

import (
	"context"
	"fmt"
	"task-manager/internal/domain"
	castomErr "task-manager/pkg/errors"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, userID uuid.UUID, req *domain.CreateTaskRequest) (*domain.Task, error)
	GetByID(ctx context.Context, userID, taskID uuid.UUID) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Task, int, error)
	Update(ctx context.Context, userID, taskID uuid.UUID, req *domain.UpdateTaskRequest) (*domain.Task, error)
	Delete(ctx context.Context, userID, taskID uuid.UUID) error
}

type TaskUsecase struct {
	repo TaskRepository
}

func NewTaskUsecase(repo TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

func (t *TaskUsecase) CreateTask(ctx context.Context, userID uuid.UUID, req *domain.CreateTaskRequest) (*domain.Task, error) {
	if req.Title == "" {
		return nil, castomErr.ErrEmptyTitle
	}
	task, err := t.repo.Create(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("task_usecase.CreateTask: %w", err)
	}
	return task, nil
}

func (t *TaskUsecase) GetTasks(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Task, int, error) {
	tasks, total, err := t.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("task_usecase.GetTasks: %w", err)
	}
	return tasks, total, nil
}

func (t *TaskUsecase) GetTaskByID(ctx context.Context, userID, taskID uuid.UUID) (*domain.Task, error) {
	task, err := t.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("task_usecase.GetTaskById: %w", err)
	}
	return task, nil
}

func (t *TaskUsecase) UpdateTask(ctx context.Context, userID, taskID uuid.UUID, req domain.UpdateTaskRequest) (*domain.Task, error) {
	if req.Title == "" {
		return nil, castomErr.ErrEmptyTitle
	}
	task, err := t.repo.Update(ctx, userID, taskID, &req)
	if err != nil {
		return nil, fmt.Errorf("task_usecase.UpdateTasc: %w", err)
	}
	return task, nil
}

func (t *TaskUsecase) DeleteTask(ctx context.Context, userID, taskID uuid.UUID) error {
	err := t.repo.Delete(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("task_usecase.DeleteTask: %w", err)
	}
	return nil
}
