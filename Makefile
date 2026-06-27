.PHONY: up down dev start swag prod

up:
	docker compose up -d

down:
	docker compose down

dev:
	docker compose up -d
	@echo "Start the app: make start"

start:
	@echo "Kill existing process on :8080..."
	-lsof -ti :8080 | xargs kill -9 2>/dev/null
	@echo "Starting app..."
	go run ./cmd/server

prod:
	docker compose --profile prod up -d

swag:
	swag init -g cmd/server/main.go --output docs
