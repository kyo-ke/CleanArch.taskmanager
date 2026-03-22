package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"taskmanager/internal/entity"
	"taskmanager/internal/infra/memory"
	"taskmanager/internal/usecase"
)

func TestTaskUsecase_CRUDAndListByUser(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewTaskRepository()
	uc := usecase.NewTaskUsecase(repo)

	created, err := uc.Create(ctx, entity.Task{
		Priority:       1,
		AssignedUserID: "user-1",
		Status:         entity.TaskStatusTodo,
		Description:    "first",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected ID to be set")
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatalf("expected CreatedAt/UpdatedAt to be set")
	}

	gotID, err := uuid.Parse(created.ID)
	if err != nil {
		t.Fatalf("created.ID should be UUID: %v", err)
	}

	got, err := uc.Get(ctx, gotID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Description != "first" {
		t.Fatalf("expected description %q, got %q", "first", got.Description)
	}

	updated, err := uc.Update(ctx, entity.Task{
		ID:             created.ID,
		Priority:       2,
		AssignedUserID: "user-2",
		Status:         entity.TaskStatusInProgress,
		Description:    "second",
		CreatedAt:      created.CreatedAt,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.AssignedUserID != "user-2" {
		t.Fatalf("expected assigned user to change")
	}

	listUser1, err := uc.ListByUser(ctx, "user-1")
	if err != nil {
		t.Fatalf("ListByUser returned error: %v", err)
	}
	if len(listUser1) != 0 {
		t.Fatalf("expected 0 tasks for user-1, got %d", len(listUser1))
	}

	listUser2, err := uc.ListByUser(ctx, "user-2")
	if err != nil {
		t.Fatalf("ListByUser returned error: %v", err)
	}
	if len(listUser2) != 1 {
		t.Fatalf("expected 1 task for user-2, got %d", len(listUser2))
	}
	if listUser2[0].ID != created.ID {
		t.Fatalf("expected same task ID")
	}

	if err := uc.Delete(ctx, gotID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	_, err = uc.Get(ctx, gotID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
	if !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
