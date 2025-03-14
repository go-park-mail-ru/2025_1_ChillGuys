# Этап 1: Сборка приложения
FROM golang:1.23.3 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код в контейнер
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./cmd/app/main.go

# Этап 2: Создание финального образа
FROM alpine:3.18

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарник из этапа сборки
COPY --from=builder /app/bin/main .

# Копируем папку с картинками 
COPY --from=builder /app/media/ ./media

# Копируем .env файл
COPY --from=builder /app/.env .

# Открываем порт для приложения
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]