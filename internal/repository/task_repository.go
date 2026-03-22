package repository

import (
	"context"

	"github.com/google/uuid"

	"taskmanager/internal/entity"
)

// TaskRepository is an output port used by usecases.
// Infrastructure (e.g. memory/postgres) provides implementations.
//
// Contract:
// - Create: returns the stored Task (with ID/CreatedAt/UpdatedAt populated)
// - Update: returns the updated Task (with UpdatedAt refreshed)
// - Get/Delete: return ErrNotFound if the record doesn't exist
// - ListByUser: returns tasks where AssignedUserID == assignedUserID
//
// Implementations must be safe for concurrent use unless stated otherwise.
type TaskRepository interface {
	Create(ctx context.Context, task entity.Task) (entity.Task, error)
	Update(ctx context.Context, task entity.Task) (entity.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (entity.Task, error)
	ListByUser(ctx context.Context, assignedUserID string) ([]entity.Task, error)
}
