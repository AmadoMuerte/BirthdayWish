.PHONY: up

db:
	docker-compose up -d

gateway:
	go run apps/gateway/cmd/main.go

wishlister:
	cd apps/wishlister && go run cmd/main.go


filer:
	cd apps/filer && go run cmd/main.go
