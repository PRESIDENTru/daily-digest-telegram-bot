FROM golang:1.25.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o tg_bot ./cmd

# Собираем бинарник
RUN CGO_ENABLED=0 go build -o /bot cmd/main.go

RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

CMD ["./tg_bog"]

