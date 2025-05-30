DATABASE_URL=postgres://gwi:password@localhost:5432/favorites?sslmode=disable
DB_HOST=127.0.0.1
PG_CONTAINER=platform-go-challenge-postgres_gwi-1

.PHONY: run build test up down logs psql reset-password check-hba

run:
	@echo Starting Go server...
	@cmd /C "set DATABASE_URL=$(DATABASE_URL)&& go run ./cmd/server"

build:
	go build -o app .

test:
	@cmd /C "set DB_HOST=$(DB_HOST)&& go test ./..."

up:
	docker-compose up -d

down:
	docker-compose down -v

logs:
	docker logs -f $(PG_CONTAINER)

psql:
	docker exec -it $(PG_CONTAINER) psql -U gwi -d favorites

reset-password:
	docker exec -it $(PG_CONTAINER) bash -c "psql -U gwi -d favorites -c \"ALTER USER gwi WITH PASSWORD 'password';\""

check-hba:
	docker exec -it $(PG_CONTAINER) cat /var/lib/postgresql/data/pg_hba.conf
