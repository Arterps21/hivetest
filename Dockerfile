# Используем официальный образ Go
FROM golang:1.20

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod .
COPY go.sum .

# Устанавливаем зависимости
RUN go mod download

# Копируем остальные файлы
COPY . .

# Компилируем приложение
RUN go build -o main .

# Запускаем приложение
CMD ["./main"]
