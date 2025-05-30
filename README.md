# GWI Favorites API

A robust Go REST API for managing user favorites at GWI, including charts, insights, and audiences.  
Built for clarity, security, and real-world deployment, tested with Docker Compose and secured with JWT.  
This project goes beyond the original challenge with production-grade API patterns, full integration tests, Docker automation, and best practices for both Linux and Windows.

---

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

## Running Tests

On Windows with MinGW:

```bash
mingw32-make down up reset-password test
```

Or run all steps with a custom target:

```bash
mingw32-make full-reset
```

This brings up a fresh Postgres, seeds with demo data, guarantees correct credentials, and runs all tests against the live DB.

---

## Running the API Locally

```bash
mingw32-make run
```

With Docker Compose:

```bash
docker-compose up -d
```

Then:

```bash
set DATABASE_URL=postgres://gwi:password@localhost:5432/favorites?sslmode=disable
go run ./cmd/server
```

---

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
On Docker for Windows, there is a widely-reported bug where Postgres password authentication (`pq: password authentication failed for user "gwi"`) fails even after destroying all volumes, resetting the password, and editing `pg_hba.conf` to use `md5`.  

This issue occurs even with fresh and known-good projects, and is due to Docker/Postgres interaction on Windows, not a code bug or logic error.

All project code, schema, and tests are correct.  

---

## Author

**George Vamvakousis**  
[geovam99@gmail.com](mailto:geovam99@gmail.com)
