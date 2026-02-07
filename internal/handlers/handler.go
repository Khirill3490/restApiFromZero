package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rest_api/internal/models"
	"rest_api/internal/repository"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store repository.TaskRepository
}

func NewHandler(store repository.TaskRepository) *Handler {
	return &Handler{store: store}
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List(r.Context(), 100, 0)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title == "" {
		respondError(w, http.StatusBadRequest, "Необходимо указать заголовок задачи")
		return
	}
	if input.Description == "" {
		respondError(w, http.StatusBadRequest, "Необходимо указать описание задачи")
		return
	}

	task, err := h.store.Create(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Ошибка при создании задачи")
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var input models.UpdateTaskInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Неверный формат данных запроса")
		return
	}

	task, err := h.store.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.store.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
