# Этап 1: Сборка
FROM golang:alpine AS builder
WORKDIR /app

# Сначала копируем файлы зависимостей для эффективного кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем чистое Go-приложение без флага -mod=vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Этап 2: Запуск
FROM alpine:3.19
WORKDIR /root/
COPY --from=builder /app/main .

# ОБЯЗАТЕЛЬНО: Указываем порт для Railway (он автоматически перенаправит трафик)
EXPOSE 8080

CMD ["./main"]
