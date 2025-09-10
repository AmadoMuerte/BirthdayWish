.PHONY: up

db:
	podman compose up -d

metrics:
	podman compose -f docker-compose.monitoring.yml up -d

.PHONY: build build-auth build-gateway build-wishlister docker-up docker-down docker-logs

build:
	@echo "Building all services..."
	go build -o bin/auth ./apps/auth/cmd/
	go build -o bin/gateway ./apps/gateway/cmd/
	go build -o bin/wishlister ./apps/wishlister/cmd/

build-auth:
	go build -o bin/auth ./apps/auth/cmd/

build-gateway:
	go build -o bin/gateway ./apps/gateway/cmd/

build-wishlister:
	go build -o bin/wishlister ./apps/wishlister/cmd/

docker-up:
	podman-compose -f docker-compose.yml up -d --build

docker-down:
	podman-compose -f docker-compose.yml down

docker-logs:
	podman-compose -f docker-compose.yml logs -f

docker-logs-auth:
	podman-compose -f docker-compose.yml logs -f auth-service

docker-logs-gateway:
	podman-compose -f docker-compose.yml logs -f gateway-service

docker-logs-wishlister:
	podman-compose -f docker-compose.yml logs -f wishlister-service

docker-restart:
	podman-compose -f docker-compose.yml restart