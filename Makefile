DATABASE_URL=postgres://gwi:password@localhost:5432/favorites?sslmode=disable
DB_HOST=127.0.0.1
PG_CONTAINER=platform-go-challenge-postgres_gwi-1
APP_NAME=gwi-favorites-api

.PHONY: run build test up down logs psql reset-password check-hba

run:
	@echo Starting Go server...
	@cmd /C "set DATABASE_URL=$(DATABASE_URL)&& go run ./cmd/server"

test:
	@cmd /C "set DB_HOST=$(DB_HOST)&& go test ./..."

up:
	docker-compose up -d

build:
	docker build -t $(APP_NAME) .

down:
	docker-compose down -v

logs:
	docker logs -f $(PG_CONTAINER)

psql:
	docker exec -it $(PG_CONTAINER) psql -U gwi -d favorites

docker-run:
	docker run --rm -p 8080:8080 \
		-e DATABASE_URL=postgres://gwi:password@host.docker.internal:5432/favorites?sslmode=disable \
		$(APP_NAME)