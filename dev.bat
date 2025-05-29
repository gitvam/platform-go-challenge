@echo off
setlocal

set IMAGE=gwi-favorites-api
set GO_VERSION=1.24.3

if "%1"=="" goto help

if "%1"=="test" (
    echo Running Go tests using Docker...
    docker run --rm ^
        -v "%cd%":/app ^
        -w /app ^
        --network=host ^
        -e DB_HOST=localhost ^
        -e JWT_SECRET=my_super_secret ^
        golang:%GO_VERSION% go test -v ./...
    goto :eof
)

if "%1"=="build" (
    echo Building Docker image...
    docker build -t %IMAGE% .
    goto :eof
)

if "%1"=="run" (
    echo Starting Docker Compose stack and running app...
    docker compose down -v
    docker compose up -d
    echo Starting Go server...
    go run cmd/server/main.go
    goto :eof
)

if "%1"=="up" (
    echo Starting Docker Compose (DB only)...
    docker compose up -d
    goto :eof
)

if "%1"=="down" (
    echo Shutting down Docker Compose and removing volumes...
    docker compose down -v
    goto :eof
)

if "%1"=="psql" (
    echo Connecting to Postgres...
    docker exec -it postgres_gwi psql -U gwi -d favorites
    goto :eof
)

:help
echo.
echo Usage: dev.bat [command]
echo ------------------------
echo test    - Run Go tests in Docker
echo build   - Build app Docker image
echo run     - Recreate DB and run the app
echo up      - Start just the DB via docker-compose
echo down    - Stop and clean docker-compose and volumes
echo psql    - Open Postgres shell inside the container
echo.
exit /b 1
