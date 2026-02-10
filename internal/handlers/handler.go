package handlers

import (
	"rest_api/internal/authjwt"
	"rest_api/internal/repository"
)

type Handler struct {
	taskStore repository.TaskRepository
	userStore repository.UserRepository
	jwt       *authjwt.Service
}

func NewHandler(taskStore repository.TaskRepository, userStore repository.UserRepository, jwtService *authjwt.Service) *Handler {
	return &Handler{
		taskStore: taskStore,
		userStore: userStore,
		jwt:       jwtService,
	}
}
