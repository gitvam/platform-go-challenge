# Build Go binary using an official lightweight image
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy go mod and sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the application
RUN go build -o server ./cmd/server/main.go

# Use a minimal runtime image
FROM alpine:3.18

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/server .

# Expose API port
EXPOSE 8080

# Start the server
CMD ["./server"]
