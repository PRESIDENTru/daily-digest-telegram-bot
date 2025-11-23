FROM golang:1.25.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Убираем дублирующую сборку, оставляем одну
RUN go build -o tg_bot ./cmd

# Исправляем опечатку в CMD
CMD ["./tg_bot"]
