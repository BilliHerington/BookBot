FROM golang:1.22.3 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
COPY .env .env

RUN go mod tidy


# Копируем весь код
COPY backend/ ./

RUN go build -o main .

CMD ["/app/main"]
