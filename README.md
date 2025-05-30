# GWI Favorites API

A robust Go REST API for managing user favorites at GWI, including charts, insights, and audiences.  
Built for clarity, security, and real-world deployment, tested with Docker Compose and secured with JWT.  
This project goes beyond the original challenge with production-grade API patterns, full integration tests, Docker automation, and best practices for both Linux and Windows.

---

## API Endpoints

All endpoints require a valid JWT token in the `Authorization: Bearer <token>` header.

| Method | Path                                              | Description                             |
|--------|---------------------------------------------------|-----------------------------------------|
| GET    | `/v1/users/{userID}/favorites`                    | List all favorite assets for the user   |
| POST   | `/v1/users/{userID}/favorites`                    | Add a new favorite asset                |
| DELETE | `/v1/users/{userID}/favorites/{assetID}?type=...` | Remove a favorite by external ID & type |
| PATCH  | `/v1/users/{userID}/favorites/{assetID}?type=...` | Edit description of a favorite asset    |

**Query Parameters:**

- `limit` / `offset` on `GET /favorites` for pagination  
- `type` on `DELETE` and `PATCH` (must be one of `chart`, `insight`, or `audience`)

> ⏳ All endpoints are protected by IP-based rate limiting: **10 requests per minute per IP**

## Project Structure

```
cmd/
  server/              
    main.go
docs/
  docs.go
  swagger.json
  swagger.yaml
handlers_test/
  api_test.go           
internal/
  handlers/
    handlers.go
  middleware/
    jwt.go
    logging.go
  models/
    asset.go
    utils.go
  store/
    postgres_store.go
    store.go
    store_test.go
  utils/
    http.go
    utils.go
Dockerfile
docker-compose.yml
.dockerignore
init.sql
Makefile
go.mod
go.sum
README.md
```

---

## Features

- REST API to add, list, edit, and remove user favorites (charts, insights, audiences)
- PostgreSQL backend with schema and seed data via `init.sql`
- Polymorphic asset model using Go interfaces
- JWT authentication with Bearer tokens in the `Authorization` header
- Swagger UI at `/swagger/index.html`
- IP-based rate limiting via `go-chi/httprate`
- Automated integration/unit tests using Dockerized Postgres and Go's `testing` package
- Consistent JSON error and success responses
- Cross-platform task automation with Makefile (works with `mingw32-make`)

---

## Best Practices Followed

- All secrets and configuration from environment variables
- Fully automated Docker Compose setup (including volume cleanup)
- Database seeding and idempotent schema via `init.sql`
- Secure JWT authentication, signature, and claims validation (middleware enforces Bearer scheme)
- Modular/idiomatic Go structure, no global state
- Live API documentation via Swagger
- Rate limiting to prevent abuse
- Tests run against a real Postgres database, not just mocks
- Scripts and Makefile work cross-platform (Windows/Linux/Mac)
- Wrapped DB initialization in transactions for safe schema creation and rollback on failure

---
## 🔍 Running Tests

### ✅ Option 1: Native Go + Docker Compose for PostgreSQL

Run tests locally with your Go toolchain, connecting to Dockerized Postgres:

```bash
mingw32-make down
mingw32-make up
mingw32-make test
```

This will:

- Reset the environment
- Start a fresh PostgreSQL container
- Run all Go tests using the local Go installation

---

## 🚀 Running the API

### 🐳 Option 1: Run API in Docker

```bash
mingw32-make build
mingw32-make docker-run
```

This will:

- Build the Go binary in a clean Docker container
- Copy Swagger docs into the image
- Run the API connected to Postgres via `host.docker.internal`
- Serve the API at `http://localhost:8080`

---

### 💻 Option 2: Run API Locally with Dockerized DB

```bash
docker-compose up -d
```

Then run the server using Go:

```bash
set DATABASE_URL=postgres://gwi:password@localhost:5432/favorites?sslmode=disable
go run ./cmd/server
```

> Swagger UI: http://localhost:8080/swagger/index.html

## JWT Authentication

The API **expects an HTTP header** with a valid Bearer token:

```
Authorization: Bearer <token>
```

- JWT secret is `my_super_secret` (demo only, use an environment variable in production).
- The `sub` claim in the JWT maps to the `userID` used for the API calls.
- If the header is missing, malformed, or the token is invalid, the API responds with `401 Unauthorized`.

### Example Payload

```json
{
  "sub": "11111111-1111-1111-1111-111111111111",
  "exp": 1999999999
}
```

You can generate test tokens at [jwt.io](https://jwt.io) using the secret `my_super_secret`.

---

## Example API Response

### Success:

```json
{
  "status": "success",
  "data": [
    { "id": "chart_engagement_2024", "title": "Q1 2024 Social Media Engagement", "type": "chart" }
  ]
}
```

### Error:

```json
{
  "status": "error",
  "message": "asset already in favorites"
}
```

---

## Additions Beyond the Original Challenge

- Automated PostgreSQL setup via Docker Compose and seed file  
- Integration tests running against a real DB  
- Rate limiting middleware for security  
- Live Swagger/OpenAPI documentation  
- Polymorphic Go types for flexible asset support  
- Consistent JSON API responses  
- Idempotent reset/test targets for easy CI/dev runs

## Future Improvements

- Add Redis or in-memory caching to reduce repeated asset lookups
- Introduce CI/CD pipeline (e.g., GitHub Actions) for automated tests, formatting, and container builds

---

## Known Issue: Docker + Windows + Postgres

**NOTE:**  
On Docker for Windows, there is a widely-reported bug where Postgres password authentication (`pq: password authentication failed for user "gwi"`) fails even after destroying all volumes, resetting the password, and editing `pg_hba.conf` to use `scram-sha-256`.

This issue occurs even with fresh and known-good projects, and is due to Docker/Postgres interaction on Windows, not a code bug or logic error.

---

### ✅ Temporary Fix (for local development only)

1. Locate the `pg_hba.conf` file inside your container or Windows installation:  
   - For Docker: use `docker exec -it <container> bash`  
   - For native installs: usually found in `C:\Program Files\PostgreSQL\XX\data`

2. Find the lines:

    ```conf
    host    all             all             127.0.0.1/32            scram-sha-256
    host    all             all             ::1/128                 scram-sha-256
    ```

3. Then change them to:

    ```conf
    host    all             all             127.0.0.1/32            trust
    host    all             all             ::1/128                 trust
    ```

4. Restart PostgreSQL:

    ```conf
    For Docker: docker restart <container-name>
    For Windows native install: run services.msc → find PostgreSQL → right-click → Restart
    ```

---

All project code, schema, and tests are correct. This repo will work out-of-the-box on Linux, macOS, WSL2, or with native Postgres.


## Author

**George Vamvakousis**  
[geovam99@gmail.com](mailto:geovam99@gmail.com)
