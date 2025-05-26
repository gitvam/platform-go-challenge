# GWI Favorites API

A Go REST API for managing user favorites at GWI, including charts, insights, and audiences. Built with clean, idiomatic Go, tested in Docker, and secured with JWT authentication.

---

## 📦 Project Structure

```
├── cmd/server           # Entry point (main.go)
├── docs/                # Swagger-generated OpenAPI docs (swagger.yaml, swagger.json, docs.go)
├── internal/
│   ├── handlers         # HTTP handlers
│   ├── middleware       # JWT auth middleware
│   ├── models           # Asset models and interface
│   ├── store            # In-memory store + seeding
│   └── utils            # JSON response wrappers, logging
├── handlers_test/       # Integration tests for the API
├── dev.bat              # Windows dev script (build, run, test)
├── Dockerfile           # Multi-stage Docker build (build + distroless run)
├── docker-compose.yml   # Service orchestration (optional)
├── go.mod               # Go module definition
├── go.sum               # Go dependency lock file
└── README.md            # Project documentation
```

---

## 🚀 Features

- ✅ Add/remove/edit favorite assets per user
- ✅ Supports Charts, Insights, Audiences via polymorphic `Asset` interface
- ✅ JWT authentication (with real signature validation)
- ✅ Swagger docs (`/swagger/index.html`)
- ✅ In-memory store with dummy data for dev
- ✅ Standardized JSON API responses
- ✅ All tests run via Docker (`go test ./...`)
- ✅ Windows-first dev flow via `dev.bat`
- 🛡️ IP-based rate limiting via `go-chi/httprate`

---

## 🧪 Running Tests

```bash
dev.bat test
```

This runs `go test -v ./...` in a Docker container using Go `1.24.3`.

---

## 🛠 Running the API

```bash
dev.bat build
dev.bat run
```

Or use Docker directly:

```bash
docker build -t gwi-favorites-api .
docker run -p 8080:8080 -e APP_ENV=dev gwi-favorites-api
```

> JWTs will be printed in the terminal on startup in dev mode.

---

## 🔐 JWT Authentication

- The API expects a valid Bearer token in `Authorization` header.
- Tokens are validated and parsed using `github.com/golang-jwt/jwt/v5`
- Subject claim (`sub`) is used as the authenticated `userID`.

### Example Token Payload

```json
{
  "sub": "johnsmith",
  "exp": 1748374243
}
```

---

## 🧾 Example API Response

### ✅ Success

```json
{
  "status": "success",
  "data": [
    { "id": "chart_1", "title": "Engagement Q1", "type": "chart" }
  ]
}
```

### ❌ Error

```json
{
  "status": "error",
  "message": "asset already in favorites"
}
```

---

## 🧠 Design Notes

- ✅ All handlers use `utils.SuccessResponse` and `utils.ErrorResponse`
- ✅ `getUserIDOrAbort()` ensures user is in context from JWT
- ✅ Assets are dynamically deserialized from JSON via a type field

---

## 💡 Future Improvements

- 🔁 Move from in-memory to persistent database (PostgreSQL)
- 🧠 Add caching layer (e.g., `go-cache` or Redis) for frequently accessed assets
- 📄 Improve Swagger schema with oneOf + discriminator

---

## 🧠 Author

George Vamvakousis  
📧 [geovam99@gmail.com](mailto:geovam99@gmail.com)
