# Генерация кода из OpenAPI спецификации
gen:
	oapi-codegen --config internal/api/oapi-codegen.yaml internal/api/openapi.yaml

# Сборка приложения
build:
	go build -o bin/gateway cmd/main.go

# Запуск сервера
run:
	go run cmd/main.go

# Установка зависимостей
deps:
	go mod tidy
	go mod download