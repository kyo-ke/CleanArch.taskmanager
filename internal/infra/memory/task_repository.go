package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"taskmanager/internal/entity"
	"taskmanager/internal/repository"
)

// TaskRepository is an in-memory implementation of repository.TaskRepository.
// Intended for local dev and unit tests.
type TaskRepository struct {
	mu    sync.RWMutex
	byID  map[uuid.UUID]entity.Task
	index map[string]map[uuid.UUID]struct{} // assignedUserID -> set(taskID)
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		byID:  make(map[uuid.UUID]entity.Task),
		index: make(map[string]map[uuid.UUID]struct{}),
	}
}

func (r *TaskRepository) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	_ = ctx

	now := time.Now().UTC()
	id, err := parseOrNewID(task.ID)
	if err != nil {
		return entity.Task{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byID[id]; exists {
		// treat as conflict; for now return existing entity as-is
		return entity.Task{}, errorsNew("task already exists")
	}

	task.ID = id.String()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}
	task.UpdatedAt = now

	r.byID[id] = task
	r.addToIndex(task.AssignedUserID, id)

	return task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task entity.Task) (entity.Task, error) {
	_ = ctx

	id, err := parseOrNewID(task.ID)
	if err != nil {
		return entity.Task{}, err
	}

	now := time.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.byID[id]
	if !ok {
		return entity.Task{}, repository.ErrNotFound
	}

	// keep CreatedAt if caller didn't set it
	if task.CreatedAt.IsZero() {
		task.CreatedAt = existing.CreatedAt
	}

	// update secondary index if user changed
	if existing.AssignedUserID != task.AssignedUserID {
		r.removeFromIndex(existing.AssignedUserID, id)
		r.addToIndex(task.AssignedUserID, id)
	}

	task.ID = id.String()
	task.UpdatedAt = now

	r.byID[id] = task
	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.byID[id]
	if !ok {
		return repository.ErrNotFound
	}

	delete(r.byID, id)
	r.removeFromIndex(existing.AssignedUserID, id)
	return nil
}

func (r *TaskRepository) Get(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.byID[id]
	if !ok {
		return entity.Task{}, repository.ErrNotFound
	}
	return task, nil
}

func (r *TaskRepository) ListByUser(ctx context.Context, assignedUserID string) ([]entity.Task, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]uuid.UUID, 0)
	if set, ok := r.index[assignedUserID]; ok {
		for id := range set {
			ids = append(ids, id)
		}
	}

	// stable-ish order: UpdatedAt desc, then ID
	tasks := make([]entity.Task, 0, len(ids))
	for _, id := range ids {
		if t, ok := r.byID[id]; ok {
			tasks = append(tasks, t)
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].UpdatedAt.Equal(tasks[j].UpdatedAt) {
			return tasks[i].ID < tasks[j].ID
		}
		return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
	})
	return tasks, nil
}

func (r *TaskRepository) addToIndex(userID string, id uuid.UUID) {
	set, ok := r.index[userID]
	if !ok {
		set = make(map[uuid.UUID]struct{})
		r.index[userID] = set
	}
	set[id] = struct{}{}
}

func (r *TaskRepository) removeFromIndex(userID string, id uuid.UUID) {
	set, ok := r.index[userID]
	if !ok {
		return
	}
	delete(set, id)
	if len(set) == 0 {
		delete(r.index, userID)
	}
}

func parseOrNewID(raw string) (uuid.UUID, error) {
	if raw == "" {
		return uuid.New(), nil
	}
	return uuid.Parse(raw)
}

// errorsNew exists to avoid importing "errors" with a name collision in this file.
func errorsNew(msg string) error { return &simpleError{s: msg} }

type simpleError struct{ s string }

func (e *simpleError) Error() string { return e.s }
