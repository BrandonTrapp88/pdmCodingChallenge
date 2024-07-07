# Define variables
API_DIR := api
FRONTEND_DIR := frontend
API_BINARY := server
FRONTEND_BUILD := build

# Check if DB_USER and DB_PASSWORD are set
ifndef DB_USER
$(error DB_USER is not set)
endif

ifndef DB_PASSWORD
$(error DB_PASSWORD is not set)
endif

# API targets
api: ## Run the API server
	@echo "Starting API server with DB_USER=$(DB_USER) and DB_PASSWORD=$(DB_PASSWORD)..."
	@cd $(API_DIR) && DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) go run main.go repository.go handlers.go routers.go

api-build: ## Build the API server
	@echo "Building API server..."
	@cd $(API_DIR) && go build -o $(API_BINARY)

api-clean: ## Clean API build artifacts
	@echo "Cleaning API build artifacts..."
	@cd $(API_DIR) && rm -f $(API_BINARY)

# Frontend targets
frontend: ## Start the frontend development server
	@echo "Starting frontend development server..."
	@cd $(FRONTEND_DIR) && npm start

frontend-build: ## Build the frontend for production
	@echo "Building frontend for production..."
	@cd $(FRONTEND_DIR) && npm run build

frontend-clean: ## Clean frontend build artifacts
	@echo "Cleaning frontend build artifacts..."
	@cd $(FRONTEND_DIR) && rm -rf $(FRONTEND_BUILD)

# General targets
all: ## Run both API and frontend
	@$(MAKE) api DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) &
	@$(MAKE) frontend &

build: api-build frontend-build ## Build both API and frontend

clean: api-clean frontend-clean ## Clean both API and frontend

help: ## Display this help message
	@echo "Usage: make [target] DB_USER=<your_db_user> DB_PASSWORD=<your_db_password>"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-20s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.PHONY: api api-build api-clean frontend frontend-build frontend-clean all build clean help
