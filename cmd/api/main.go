package main

func main() {
	// This is a placeholder for the main function of the API command.
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

