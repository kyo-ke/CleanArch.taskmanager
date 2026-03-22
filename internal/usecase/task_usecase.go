package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"taskmanager/internal/entity"
	"taskmanager/internal/repository"
)

// TaskUsecase provides application operations for Tasks.
// It depends only on the repository port (interface) and entities.
type TaskUsecase struct {
	repo repository.TaskRepository
	now  func() time.Time
}

func NewTaskUsecase(repo repository.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: repo, now: func() time.Time { return time.Now().UTC() }}
}

func (u *TaskUsecase) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	if task.Status == "" {
		task.Status = entity.TaskStatusTodo
	}
	if err := task.Validate(); err != nil {
		return entity.Task{}, err
	}
	now := u.now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}
	task.UpdatedAt = now
	return u.repo.Create(ctx, task)
}

func (u *TaskUsecase) Update(ctx context.Context, task entity.Task) (entity.Task, error) {
	if err := task.Validate(); err != nil {
		return entity.Task{}, err
	}
	task.UpdatedAt = u.now()
	updated, err := u.repo.Update(ctx, task)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return entity.Task{}, ErrNotFound
		}
		return entity.Task{}, err
	}
	return updated, nil
}

func (u *TaskUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (u *TaskUsecase) Get(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	task, err := u.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return entity.Task{}, ErrNotFound
		}
		return entity.Task{}, err
	}
	return task, nil
}

func (u *TaskUsecase) ListByUser(ctx context.Context, assignedUserID string) ([]entity.Task, error) {
	return u.repo.ListByUser(ctx, assignedUserID)
}
