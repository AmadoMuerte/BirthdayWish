.PHONY: up

db:
	podman compose up -d

metrics:
	podman compose -f docker-compose.monitoring.yml up -d

.PHONY: docker-up docker-down docker-logs

run-auth:
	cd apps/auth && go run cmd/main.go

run-gateway:
	cd apps/gateway && go run cmd/main.go

run-wishlister:
	cd apps/wishlister && go run cmd/main.go

docker-up:
	podman-compose -f docker-compose.yml up --build

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