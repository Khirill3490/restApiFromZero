package handlers

import (
	"rest_api/internal/authjwt"
	"rest_api/internal/repository"
)

type Handler struct {
	taskStore    repository.TaskRepository
	userStore    repository.UserRepository
	refreshStore repository.RefreshSessionRepository
	jwt          *authjwt.Service
}

func NewHandler(
	taskStore repository.TaskRepository,
	userStore repository.UserRepository,
	refreshStore repository.RefreshSessionRepository,
	jwt *authjwt.Service,
) *Handler {
	return &Handler{
		taskStore:    taskStore,
		userStore:    userStore,
		refreshStore: refreshStore,
		jwt:          jwt,
	}
}
