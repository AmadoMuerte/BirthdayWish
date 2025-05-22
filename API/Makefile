.PHONY: up

db:
	docker-compose up -d

run-getaway:
	go run cmd/main.go

run-srvc-wishlist:
	cd services/wishlist && go run cmd/main.go


run-srvc-minio:
	cd services/minio && go run cmd/main.go
