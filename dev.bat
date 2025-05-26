@echo off
setlocal

set IMAGE=gwi-favorites-api
set GO_VERSION=1.24.3

REM Help screen
if "%1"=="" goto help

if "%1"=="test" (
    echo Running Go tests using Docker...
    docker run --rm -v "%cd%":/app -w /app golang:%GO_VERSION% go test -v ./...
    goto :eof
)

if "%1"=="build" (
    echo Building Docker image...
    docker build -t %IMAGE% .
    goto :eof
)

if "%1"=="run" (
    echo Running app container...
    docker run --rm --name %IMAGE% -e APP_ENV=dev -p 8080:8080 %IMAGE%
    goto :eof
)

if "%1"=="stop" (
    echo Stopping and removing container...
    docker stop %IMAGE%
    docker rm %IMAGE%
    goto :eof
)

:help
echo.
echo Usage: dev.bat [command]
echo ------------------------
echo test   - Run Go tests in Docker
echo build  - Build Docker image
echo run    - Run app with port 8080 and APP_ENV=dev
echo stop   - Stop and remove the container
echo.
exit /b 1