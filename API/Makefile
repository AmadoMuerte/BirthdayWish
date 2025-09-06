.PHONY: up

db:
	podman compose up -d

# # need to install parallel: 
# # brew install parallel 
# # sudo apt install parallel
# # npm install -g parallel
# # pacman -S parallel 
# run-all:
# 	@parallel --linebuffer --tagstring '{#}' ::: \
# 		"cd apps/gateway && go run cmd/main.go" \
# 		"cd apps/wishlister && go run cmd/main.go" \
# 		"cd apps/filer && go run cmd/main.go"


run-gateway:
	cd apps/gateway && go build -o ./bin/gateway ./cmd/main.go && ./bin/gateway
