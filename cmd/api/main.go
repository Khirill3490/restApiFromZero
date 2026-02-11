package main

import (
	"log"
	"net/http"

	"rest_api/internal/authjwt"
	"rest_api/internal/config"
	"rest_api/internal/db"
	"rest_api/internal/db/sqlc"
	"rest_api/internal/handlers"
	"rest_api/internal/httpserver"
)

func main() {
	cfg := config.MustLoad()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer conn.Close()

	queries := sqlc.New(conn)

	jwtSvc := authjwt.New(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	taskStore := db.NewTaskStore(queries)
	userStore := db.NewUserStore(queries)

	h := handlers.NewHandler(taskStore, userStore, jwtSvc)

	router := httpserver.NewRouter(h)

	log.Println("Сервер запущен на адресе", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, router))
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
