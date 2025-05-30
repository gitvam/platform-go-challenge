FROM golang:1.24.3-alpine AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/server/main.go -o docs

RUN go build -o app ./cmd/server

FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/docs ./docs

EXPOSE 8080
CMD ["/app/app"]
