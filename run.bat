@echo off

if "%1"=="" goto help

if "%1"=="dev" goto dev
if "%1"=="start" goto start
if "%1"=="stop" goto stop
if "%1"=="swag" goto swag

echo Unknown command: %1
:help
echo Usage: run.bat {dev^|start^|stop^|swag}
echo   dev   Start infrastructure ^(MySQL + Redis + RabbitMQ^)
echo   start Start app ^(auto-clear port 8080^)
echo   stop  Stop infrastructure
echo   swag  Generate Swagger docs
exit /b 1

:dev
docker exec eshop-mysql mysqladmin ping -uroot -p123456 --silent >nul 2>&1
if not errorlevel 1 (
    echo Infrastructure already running.
    echo Start the app: run.bat start
    exit /b 0
)
echo [1/2] Shutdown WSL to free ports...
wsl --terminate Ubuntu-24.04
echo [2/2] Start infrastructure...
docker compose up -d
echo.
echo Done. Start the app: run.bat start
exit /b 0

:start
echo [1/2] Kill existing process on port 8080...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":8080"') do (
    taskkill /F /PID %%a >nul 2>&1
)
echo [2/2] Starting app...
go run ./cmd/server
exit /b 0

:stop
echo Stopping infrastructure...
docker compose down
echo Done.
exit /b 0

:swag
swag init -g cmd/server/main.go --output docs
exit /b 0
