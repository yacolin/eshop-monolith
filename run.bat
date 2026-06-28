@echo off

if "%1"=="" goto help

if "%1"=="start" goto start
if "%1"=="stop" goto stop
if "%1"=="swag" goto swag

echo Unknown command: %1
:help
echo Usage: run.bat {start^|stop^|swag}
echo   start Start app ^(auto-clear port 8080^)
echo   stop  Stop app on port 8080
echo   swag  Generate Swagger docs
exit /b 1


:start
echo [1/2] Kill existing process on port 8080...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":8080"') do (
    taskkill /F /PID %%a >nul 2>&1
)
echo [2/2] Starting app...
go run ./cmd/server
exit /b 0


:stop
echo Stopping app on port 8080...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":8080"') do (
    taskkill /F /PID %%a >nul 2>&1
)
echo Stopped.
exit /b 0


:swag
swag init -g cmd/server/main.go --output docs
exit /b 0
