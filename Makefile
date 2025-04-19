appName = lms
compose = docker-compose -f docker-compose.yml -p $(appName)

up: down build
	@echo "Starting app..."
	$(compose) up -d
	@echo "Docker images built and started!"

build:
	@echo "Building images"
	$(compose) build
	@echo "Docker images built!"

down:
	@echo "Stopping docker compose..."
	$(compose) down
	@echo "Done!"

down-v:
	@echo "Stopping docker compose and removing volumes..."
	$(compose) down -v
	@echo "Done!"