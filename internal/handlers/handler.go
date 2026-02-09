package handlers

import (
	"rest_api/internal/repository"
)

type Handler struct {
	taskStore repository.TaskRepository
	userStore repository.UserRepository
}

func NewHandler(taskStore repository.TaskRepository, userStore repository.UserRepository) *Handler {
	return &Handler{
		taskStore: taskStore,
		userStore: userStore,
	}
}
