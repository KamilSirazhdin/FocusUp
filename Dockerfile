FROM golang:1.23-alpine AS builder

WORKDIR /app

# Установка необходимых пакетов
RUN apk add --no-cache gcc musl-dev git

# Копирование и загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Установка Goose для миграций
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.17.0

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o main ./cmd/api

# Многоэтапная сборка для уменьшения размера образа
FROM alpine:latest

WORKDIR /app

# Установка необходимых зависимостей для работы приложения
RUN apk add --no-cache ca-certificates tzdata

# Копирование исполняемого файла из предыдущего этапа
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копирование миграций и конфигураций
COPY migrations ./migrations
COPY .env.example .env

# Создание непривилегированного пользователя
RUN adduser -D appuser
USER appuser

# Запуск приложения
CMD ["./main"]