# Этап 1: Сборка приложения
FROM golang:1.23.3 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код в контейнер
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем основное приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./cmd/app/main.go

# Собираем приложение для миграций
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate ./cmd/migrations/main.go

# Этап 2: Создание финального образа
FROM alpine:3.18

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранные бинарники из этапа сборки
COPY --from=builder /app/bin/main .
COPY --from=builder /app/bin/migrate .

# Копируем папку с миграциями
COPY --from=builder /app/db/migrations ./db/migrations

# Копируем .env файл
COPY --from=builder /app/.env .

# Открываем порт для приложения
EXPOSE 8080