.PHONY: up

db:
	docker-compose up -d

run-gateway:
	go run cmd/main.go

run-srvc-wishlist:
	cd ms/wishlist && go run cmd/main.go


run-srvc-minio:
	cd ms/minio && go run cmd/main.go
