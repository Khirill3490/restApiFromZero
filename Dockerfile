# ===== build stage =====
FROM golang:1.25.4-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /app/bin/task-api ./cmd/api

# ===== runtime stage =====
FROM alpine:3.22

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/task-api /app/task-api

EXPOSE 8080
ENTRYPOINT ["/app/task-api"]
