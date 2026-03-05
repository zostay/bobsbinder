.PHONY: up down restart logs logs-api logs-frontend logs-marketing migrate-up migrate-down migrate-new seed test lint clean help

up: ## Start all services
	docker compose up -d --build

down: ## Stop all services
	docker compose down

restart: down up ## Restart all services

logs: ## Follow logs for all services
	docker compose logs -f

logs-api: ## Follow API logs
	docker compose logs -f api

logs-frontend: ## Follow frontend logs
	docker compose logs -f frontend

logs-marketing: ## Follow marketing site logs
	docker compose logs -f marketing

migrate-up: ## Run database migrations
	docker compose exec api sql-migrate up -config=/app/dbconfig.yaml

migrate-down: ## Rollback last migration
	docker compose exec api sql-migrate down -config=/app/dbconfig.yaml -limit=1

migrate-new: ## Create a new migration file (usage: make migrate-new NAME=my-migration)
	@test -n "$(NAME)" || (echo "Usage: make migrate-new NAME=my-migration" && exit 1)
	@echo "-- +migrate Up\n\n-- +migrate Down" > migrations/$$(date +%Y%m%d%H%M%S)-$(NAME).sql
	@echo "Created migrations/$$(date +%Y%m%d%H%M%S)-$(NAME).sql"

seed: ## Seed database with sample data
	docker compose exec api sql-migrate up -config=/app/dbconfig.yaml

test: ## Run Go tests
	go test ./...

lint: ## Run Go linter
	golangci-lint run ./...

clean: ## Remove containers, volumes, and temp files
	docker compose down -v
	rm -rf tmp/
	rm -rf frontend/node_modules frontend/dist
	rm -rf marketing/public marketing/resources

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
