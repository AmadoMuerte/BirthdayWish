.PHONY: up

db:
	podman compose up -d

metrics:
	podman compose -f docker-compose.monitoring.yml up -d

.PHONY: docker-up docker-down docker-logs

run-all:
	concurrently --prefix-colors "blue,green,red,yellow" \
  --names "GATEWAY,WISHLISTER,FILER,AUTH" \
  "cd apps/gateway && go run cmd/main.go" \
  "cd apps/wishlister && go run cmd/main.go" \
  "cd apps/filer && go run cmd/main.go" \
  "cd apps/auth && go run cmd/main.go"


run-filer:
	cd apps/filer && go run cmd/main.go

run-auth:
	cd apps/auth && go run cmd/main.go

run-gateway:
	cd apps/gateway && go run cmd/main.go

run-wishlister:
	cd apps/wishlister && go run cmd/main.go

docker-db:
	podman-compose -f docker-compose-db.yml up -d

docker-up:
	podman-compose -f docker-compose.yml up --build