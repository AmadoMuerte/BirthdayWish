# Практическое руководство по Docker Compose для BirthdayWish

## Пошаговая инструкция для новичков

### Шаг 1: Установка Docker

**На Ubuntu/Debian:**
```bash
# Обновление пакетов
sudo apt update

# Установка зависимостей
sudo apt install apt-transport-https ca-certificates curl gnupg lsb-release

# Добавление GPG ключа Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Добавление репозитория Docker
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Установка Docker
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Добавление пользователя в группу docker
sudo usermod -aG docker $USER
# Перезайдите в систему для применения изменений
```

**Проверка установки:**
```bash
docker --version
docker compose version
```

### Шаг 2: Понимание структуры проекта

Ваш проект использует микросервисную архитектуру:

```
BirthdayWish/
├── API/                          # Backend сервисы
│   ├── apps/                     # Микросервисы
│   │   ├── auth/                 # Сервис аутентификации
│   │   ├── wishlister/           # Сервис желаний
│   │   ├── filer/                # Сервис файлов
│   │   └── gateway/              # API Gateway
│   ├── docker-compose.yml        # Основная конфигурация
│   └── monitoring/               # Конфигурация мониторинга
└── Web/                          # Frontend приложение
```

### Шаг 3: Анализ docker-compose.yml

#### Базы данных

**PostgreSQL (auth-db и wish-db):**
```yaml
auth-db:
  image: postgres:17-alpine        # Образ PostgreSQL
  environment:                     # Переменные окружения
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: users_service
  ports:
    - "5433:5432"                  # Порт хоста:порт контейнера
  volumes:
    - auth_db_data:/var/lib/postgresql/data  # Постоянное хранение
  healthcheck:                     # Проверка готовности
    test: ["CMD-SHELL", "pg_isready -U postgres"]
    interval: 10s                  # Проверка каждые 10 секунд
    timeout: 10s                   # Таймаут 10 секунд
    retries: 10                    # 10 попыток
```

**MinIO (filer-db):**
```yaml
filer-db:
  image: minio/minio               # Образ MinIO (S3-совместимое хранилище)
  environment:
    - MINIO_ROOT_USER=minioadmin
    - MINIO_ROOT_PASSWORD=minioadmin
  ports:
    - "9000:9000"                  # API порт
    - "9001:9001"                  # Web консоль
  command: server /data --console-address ":9001"
```

#### Микросервисы

**Auth Service:**
```yaml
auth-service:
  build:
    dockerfile: apps/auth/Dockerfile  # Сборка из Dockerfile
  ports:
    - "50051:50051"               # gRPC порт
  depends_on:
    auth-db:
      condition: service_healthy   # Ждет готовности БД
```

**Gateway Service:**
```yaml
gateway-service:
  build:
    dockerfile: apps/gateway/Dockerfile
  ports:
    - "3030:3030"                 # HTTP API порт
  depends_on:
    filer-service:
      condition: service_started   # Ждет запуска всех сервисов
    auth-service:
      condition: service_started
    wishlister-service:
      condition: service_started
```

### Шаг 4: Запуск проекта

#### Первый запуск

```bash
# Переход в директорию API
cd /home/amado/Projects/BirthdayWish/API

# Запуск всех сервисов в фоновом режиме
docker compose up -d

# Просмотр статуса сервисов
docker compose ps
```

#### Проверка работы

```bash
# Просмотр логов всех сервисов
docker compose logs

# Просмотр логов конкретного сервиса
docker compose logs auth-service

# Просмотр логов в реальном времени
docker compose logs -f gateway-service
```

#### Проверка доступности

```bash
# Проверка API Gateway
curl http://localhost:3030/health

# Проверка MinIO консоли
# Откройте http://localhost:9001 в браузере
# Логин: minioadmin, Пароль: minioadmin

# Проверка Prometheus (если запущен)
# Откройте http://localhost:9090 в браузере
```

### Шаг 5: Управление сервисами

#### Остановка и запуск

```bash
# Остановка всех сервисов
docker compose down

# Остановка с удалением volumes (ОСТОРОЖНО: удалит данные!)
docker compose down -v

# Запуск только баз данных
docker compose -f docker-compose-db.yml up -d

# Запуск с мониторингом
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml up -d
```

#### Пересборка образов

```bash
# Пересборка всех образов
docker compose build

# Пересборка конкретного сервиса
docker compose build auth-service

# Пересборка без кеша
docker compose build --no-cache auth-service
```

### Шаг 6: Отладка и мониторинг

#### Просмотр ресурсов

```bash
# Использование ресурсов контейнерами
docker stats

# Детальная информация о контейнере
docker inspect auth-service

# Процессы в контейнере
docker exec -it auth-service ps aux
```

#### Подключение к контейнерам

```bash
# Подключение к контейнеру
docker exec -it auth-service sh

# Подключение к базе данных
docker exec -it auth-db psql -U postgres -d users_service
```

### Шаг 7: Настройка переменных окружения

#### Создание .env файла

```bash
# Создайте файл .env в директории API
cat > .env << EOF
# Database settings
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_password_123
POSTGRES_DB_USERS=users_service
POSTGRES_DB_WISHES=wish_service

# MinIO settings
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=secure_minio_password

# Service ports
AUTH_PORT=50051
WISHLISTER_PORT=50052
FILER_PORT=50053
GATEWAY_PORT=3030

# Database ports
AUTH_DB_PORT=5433
WISH_DB_PORT=5434
MINIO_API_PORT=9000
MINIO_CONSOLE_PORT=9001
EOF
```

#### Обновление docker-compose.yml для использования .env

```yaml
services:
  auth-db:
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB_USERS}
    ports:
      - "${AUTH_DB_PORT}:5432"
```

### Шаг 8: Добавление новых сервисов

#### Пример: Добавление Redis для кеширования

```yaml
# Добавьте в docker-compose.yml
services:
  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - birthdaywish

# Добавьте volume
volumes:
  redis_data:
    driver: local
```

### Шаг 9: Подготовка к продакшену

#### Создание docker-compose.prod.yml

```yaml
version: '3.8'

services:
  auth-service:
    build:
      dockerfile: apps/auth/Dockerfile
    environment:
      - MODE=production
      - LOG_LEVEL=info
    restart: always
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
        reservations:
          memory: 256M
          cpus: '0.25'
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

#### Использование внешних образов

```yaml
services:
  auth-service:
    image: your-registry/birthdaywish-auth:latest
    # вместо build:
```

### Шаг 10: Мониторинг и логирование

#### Настройка централизованного логирования

```yaml
services:
  loki:
    image: grafana/loki:2.9.0
    ports:
      - "3100:3100"
    volumes:
      - loki_data:/loki
    command: -config.file=/etc/loki/local-config.yaml

  promtail:
    image: grafana/promtail:2.9.0
    volumes:
      - /var/log:/var/log:ro
      - ./monitoring/promtail.yml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
```

## Частые проблемы и решения

### Проблема: Сервис не запускается

**Решение:**
```bash
# Проверьте логи
docker compose logs service-name

# Проверьте статус
docker compose ps

# Перезапустите сервис
docker compose restart service-name
```

### Проблема: Порты заняты

**Решение:**
```bash
# Найдите процесс, использующий порт
sudo netstat -tulpn | grep :3030

# Убейте процесс
sudo kill -9 PID

# Или измените порт в docker-compose.yml
```

### Проблема: Нехватка места на диске

**Решение:**
```bash
# Очистка неиспользуемых образов
docker image prune -a

# Очистка неиспользуемых volumes
docker volume prune

# Очистка всего
docker system prune -a
```

## Следующие шаги для изучения Kubernetes

1. **Установите Minikube** для локальной разработки
2. **Изучите основные ресурсы K8s**: Pod, Service, Deployment
3. **Создайте манифесты** для ваших сервисов
4. **Настройте Ingress** для внешнего доступа
5. **Добавьте мониторинг** с Prometheus Operator

Помните: практика - ключ к пониманию. Начните с простых примеров, постепенно усложняйте задачи.
