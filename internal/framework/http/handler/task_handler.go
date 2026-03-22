package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"taskmanager/internal/entity"
	"taskmanager/internal/usecase"
)

type TaskHandler struct {
	uc *usecase.TaskUsecase
}

func NewTaskHandler(uc *usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{uc: uc}
}

type taskDTO struct {
	ID             string            `json:"id,omitempty"`
	Priority       entity.Priority   `json:"priority"`
	AssignedUserID string            `json:"assignedUserId"`
	Status         entity.TaskStatus `json:"status"`
	Description    string            `json:"description"`
	CreatedAt      string            `json:"createdAt,omitempty"`
	UpdatedAt      string            `json:"updatedAt,omitempty"`
}

func fromEntity(t entity.Task) taskDTO {
	d := taskDTO{
		ID:             t.ID,
		Priority:       t.Priority,
		AssignedUserID: t.AssignedUserID,
		Status:         t.Status,
		Description:    t.Description,
	}
	if !t.CreatedAt.IsZero() {
		d.CreatedAt = t.CreatedAt.Format(timeRFC3339)
	}
	if !t.UpdatedAt.IsZero() {
		d.UpdatedAt = t.UpdatedAt.Format(timeRFC3339)
	}
	return d
}

func (d taskDTO) toEntity() entity.Task {
	return entity.Task{
		ID:             d.ID,
		Priority:       d.Priority,
		AssignedUserID: d.AssignedUserID,
		Status:         d.Status,
		Description:    d.Description,
	}
}

const timeRFC3339 = time.RFC3339

func (h *TaskHandler) Register(r chi.Router) {
	r.Post("/tasks", h.create)
	r.Put("/tasks/{id}", h.update)
	r.Delete("/tasks/{id}", h.delete)
	r.Get("/tasks/{id}", h.get)
	r.Get("/users/{userId}/tasks", h.listByUser)
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	var dto taskDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	t, err := h.uc.Create(r.Context(), dto.toEntity())
	if err != nil {
		writeJSONError(w, mapErrorToStatus(err), err)
		return
	}
	writeJSON(w, http.StatusCreated, fromEntity(t))
}

func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	var dto taskDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}
	dto.ID = id.String()

	t, err := h.uc.Update(r.Context(), dto.toEntity())
	if err != nil {
		writeJSONError(w, mapErrorToStatus(err), err)
		return
	}
	writeJSON(w, http.StatusOK, fromEntity(t))
}

func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeJSONError(w, mapErrorToStatus(err), err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	t, err := h.uc.Get(r.Context(), id)
	if err != nil {
		writeJSONError(w, mapErrorToStatus(err), err)
		return
	}
	writeJSON(w, http.StatusOK, fromEntity(t))
}

func (h *TaskHandler) listByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		writeJSONError(w, http.StatusBadRequest, errors.New("userId is required"))
		return
	}

	tasks, err := h.uc.ListByUser(r.Context(), userID)
	if err != nil {
		writeJSONError(w, mapErrorToStatus(err), err)
		return
	}

	out := make([]taskDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, fromEntity(t))
	}
	writeJSON(w, http.StatusOK, out)
}

func mapErrorToStatus(err error) int {
	if errors.Is(err, usecase.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}
