FROM golang:1.23-alpine AS builder

WORKDIR /app/back-end
# Копируем файл go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь код в контейнер
COPY . ./

# Сборка Go-приложения
RUN go build ./cmd/apiserver/.

# Создаем минимальный образ для продакшн
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированное бинарное приложение из этапа сборки
COPY --from=builder /app/back-end/apiserver .

# Экспонируем порт, который используется приложением
EXPOSE 8080

# Запускаем приложение
CMD ["./apiserver"]
