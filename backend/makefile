.PHONY: help run build migrate-up migrate-down migrate-status migrate-reset migrate-create


# Default target
help:
	@echo "Available commands:"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application"
	@echo "  make migrate-up    - Run all pending migrations"
	@echo "  make migrate-down  - Rollback last migration"

include .env
export

DB_URL := $(DATABASE_URL)

# Run application
run:
	go run cmd/api/main.go

# Build application
build:
	go build -o bin/api cmd/api/main.go

# Run migrations up
migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

# Rollback last migration
migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

# Migration status
migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

# Reset all migrations
migrate-reset:
	goose -dir migrations postgres "$(DB_URL)" reset

# Create new migration (usage: make migrate-create NAME=create_something)
migrate-create:
	goose -dir migrations create $(NAME) sql