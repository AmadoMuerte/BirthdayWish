
# Сборка приложения
build:
	go build -o bin/auth cmd/main.go

# Запуск сервера
run:
	go run cmd/main.go -mode development

# Установка зависимостей
deps:
	go mod tidy
	go mod download