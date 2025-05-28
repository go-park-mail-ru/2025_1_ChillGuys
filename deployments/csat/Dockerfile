# Этап сборки
FROM golang:1.23.3 AS builder

WORKDIR /app

# Копируем только go.mod и go.sum для оптимизации кэширования слоев
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код микросервиса user и общие внутренние пакеты
COPY . .

# Собираем микросервис
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/csat-service ./cmd/csat/main.go

# Собираем приложение для миграций
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate ./cmd/migrations/main.go

# Финальный образ
FROM alpine:3.18

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/csat-service .
COPY --from=builder /app/bin/migrate .

# Копируем папку с миграциями
COPY --from=builder /app/db/migrations ./db/migrations

# Порт gRPC
EXPOSE 50053

CMD ["./csat-service"]