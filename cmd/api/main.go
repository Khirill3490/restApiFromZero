package main

import (
	"log"
	"net/http"
	"os"
	"rest_api/internal/db"
	"rest_api/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	conn, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer conn.Close()

	store := db.NewTaskStore(conn)
	handler := handlers.NewHandler(store)

	r := chi.NewRouter()

	// Полезные middleware
	r.Use(middleware.Logger)    // логирует все запросы
	r.Use(middleware.Recoverer) // не даёт серверу упасть при panic

	// 6) Роуты
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", handler.GetAllTasks)   // GET /tasks
		r.Post("/", handler.CreateTask)  // POST /tasks

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetTaskByID) // GET /tasks/{id}
			r.Patch("/", handler.UpdateTask) // PATCH /tasks/{id}
			r.Delete("/", handler.DeleteTask) // DELETE /tasks/{id}
		})
	})

	// 7) Старт сервера
	addr := ":8080"
	log.Println("Server started on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// yourapp/
// ├── cmd/api/main.go
// ├── internal/
// │   ├── domain/
// │   │   └── user.go
// │   ├── handler/http/
// │   │   └── user.go
// │   ├── service/
// │   │   └── user.go
// │   ├── repository/
// │   │   ├── user.go        ← интерфейс репо
// │   │   └── errors.go      ← ErrNotFound и т.п.
// │   └── repositoryimpl/
// │       └── postgres/
// │           ├── user_repo.go  ← реализация репо (обёртка над sqlc)
// │           └── db/
// │               ├── sqlc/       ← (generated) сюда кладём код sqlc
// │               │   ├── db.go
// │               │   ├── models.go
// │               │   └── queries.sql.go
// │               ├── schema/
// │               │   └── schema.sql
// │               └── queries/
// │                   └── user.sql
// ├── sqlc.yaml
// └── migrations/ (если используешь миграции отдельно)

