.PHONY: help go-dev-up go-dev-up-d all-up all-up-d test build run health-check generate-report report-download logs logs-all status go-dev-down all-down rebuild

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


# Docker Compose Commands
go-dev-up: ## Start development environment (go-dev profile)
	sudo docker-compose -f docker-compose-dev.yaml --profile go-dev up --build 

go-dev-up-d: ## Start development environment in detached mode
	sudo docker-compose -f docker-compose-dev.yaml --profile go-dev up --build -d

all-up: ## Start all services (all profile)
	sudo docker-compose -f docker-compose-dev.yaml --profile all up --build 

all-up-d: ## Start all services in detached mode
	sudo docker-compose -f docker-compose-dev.yaml --profile all up --build -d

go-dev-down: ## Stop development environment
	sudo docker-compose -f docker-compose-dev.yaml --profile go-dev down

all-down: ## Stop all services and remove volumes
	sudo docker-compose -f docker-compose-dev.yaml --profile all down -v

rebuild: ## Rebuild and restart go-pdf-service only
	sudo docker-compose -f docker-compose-dev.yaml up --build go-pdf-service --profile go-dev

# Testing Commands
test: ## Run all tests
	go test ./...

test-verbose: ## Run all tests with verbose output
	go test -v ./...

# Local Development Commands
build: ## Build the Go application locally
	go build -o bin/main ./cmd

run: ## Run the application locally
	go run ./cmd

clean: ## Clean build artifacts
	rm -f bin/main
	rm -rf tmp/
	rm -rf reports/*.pdf

# Docker Commands
logs: ## Show logs from go-pdf-service
	sudo docker-compose -f docker-compose-dev.yaml logs -f go-pdf-service

logs-all: ## Show logs from all services
	sudo docker-compose -f docker-compose-dev.yaml logs -f

status: ## Show status of all containers
	sudo docker-compose -f docker-compose-dev.yaml ps

# API Testing Commands
health-check: ## Test health endpoint
	@echo "Testing health endpoint..."
	curl -s http://localhost:8080/api/v1/health

generate-report: ## Generate PDF report for student ID 1
	@echo "Generating PDF report for student ID 1..."
	curl -s http://localhost:8080/api/v1/students/1/report

report-download: ## Download PDF report for student ID 1
	@echo "Downloading PDF report for student ID 1..."
	curl -o student_1_report.pdf "http://localhost:8080/api/v1/students/1/report?download=true"
	@echo "PDF downloaded as: student_1_report.pdf"

# Quick Commands
quick-start: ## Quick start development environment
	@echo "🚀 Starting development environment..."
	@make dev-up-d
	@echo "Waiting for services..."
	@sleep 30
	@make health-check
	@echo "✅ Ready! Use 'make logs' to view logs, 'make dev-down' to stop"

quick-test: ## Quick test of PDF generation
	@make health-check
	@make generate-report 