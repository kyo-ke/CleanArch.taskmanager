package entity

import (
	"errors"
	"strings"
	"time"
)

// Priority represents how urgent/important a Task is.
// Kept as an int to stay storage/API friendly.
// Typical values: 1 (high) .. 5 (low).
type Priority int

// TaskStatus represents the current state of a Task.
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// Task is the core entity.
// Note: AssignedUserID is required for per-user task listing.
type Task struct {
	ID             string
	Priority       Priority
	AssignedUserID string
	Status         TaskStatus
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (t Task) Validate() error {
	if strings.TrimSpace(t.AssignedUserID) == "" {
		return errors.New("assigned user id is required")
	}
	if strings.TrimSpace(t.Description) == "" {
		return errors.New("description is required")
	}
	if t.Priority == 0 {
		return errors.New("priority is required")
	}
	switch t.Status {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return nil
	default:
		return errors.New("invalid status")
	}
}
