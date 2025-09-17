ifndef CONTAINER_ENGINE
    ifneq (, $(shell which podman 2>/dev/null))
        CONTAINER_ENGINE := podman
        HAS_STANDALONE_COMPOSE := $(shell which podman-compose 2>/dev/null)
        ifdef HAS_STANDALONE_COMPOSE
            COMPOSE_ENGINE := podman-compose
        else
            COMPOSE_ENGINE := podman compose
        endif
    else ifneq (, $(shell which docker 2>/dev/null))
        CONTAINER_ENGINE := docker
        HAS_DOCKER_COMPOSE := $(shell which docker compose 2>/dev/null)
        ifdef HAS_DOCKER_COMPOSE
            COMPOSE_ENGINE := docker compose
        else
            COMPOSE_ENGINE := docker compose
        endif
    else
        $(error "Ни Podman, ни Docker не найдены в системе!")
    endif
endif

ENGINE := $(CONTAINER_ENGINE)

.PHONY: help up db metrics run-all run-filer run-auth run-gateway run-wishlister

help:
	@echo "Доступные команды:"
	@echo "  make up          - Запуск всех сервисов"
	@echo "  make db          - Запуск только базы данных"
	@echo "  make metrics     - Запуск систем мониторинга"
	@echo "  make run-all     - Запуск всех Go-приложений"
	@echo "  make run-filer   - Запуск только filer"
	@echo "  make run-auth    - Запуск только auth"
	@echo "  make run-gateway - Запуск только gateway"
	@echo "  make run-wishlister - Запуск только wishlister"
	@echo ""
	@echo "Используемый контейнерный движок: $(CONTAINER_ENGINE)"
	@echo "Используемый compose движок: $(COMPOSE_ENGINE)"

up:
	$(COMPOSE_ENGINE) -f docker-compose.yml up --build

db:
	$(COMPOSE_ENGINE) -f docker-compose-db.yml up -d

metrics:
	$(COMPOSE_ENGINE) -f docker-compose.monitoring.yml up -d

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