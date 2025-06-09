.PHONY: up

db:
	docker-compose up -d

# need to install parallel: 
# brew install parallel 
# sudo apt install parallel
# npm install -g parallel
# pacman -S parallel 
run-all:
	@parallel --linebuffer --tagstring '{#}' ::: \
		"cd apps/gateway && go run cmd/main.go" \
		"cd apps/wishlister && go run cmd/main.go" \
		"cd apps/filer && go run cmd/main.go"


gateway:
	go run apps/gateway/cmd/main.go

wishlister:
	cd apps/wishlister && go run cmd/main.go


filer:
	cd apps/filer && go run cmd/main.go

