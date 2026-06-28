.PHONY: start stop swag

start:
	@echo "Starting app..."
	go run ./cmd/server

stop:
	@echo "Stopping app on :8080..."
	-@lsof -ti :8080 | xargs kill -9 2>/dev/null

swag:
	swag init -g cmd/server/main.go --output docs
