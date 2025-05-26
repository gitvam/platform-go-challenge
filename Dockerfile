# Build Go binary using an official lightweight image
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run tests – fail build if any test fails
RUN go test -v ./...

# Build the app
RUN go build -o app ./cmd/server

# Stage 2: Create a small final image
FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /app/app /app/

# Expose port 8080 for the API
EXPOSE 8080

CMD ["/app/app"]
