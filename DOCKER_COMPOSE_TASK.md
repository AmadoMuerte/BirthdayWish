# Задача: Создание и настройка Docker Compose для микросервисной архитектуры

## Обзор проекта

У вас есть микросервисное приложение BirthdayWish, состоящее из нескольких компонентов:
- **auth-service** - сервис аутентификации (gRPC на порту 50051)
- **wishlister-service** - сервис управления желаниями (gRPC на порту 50052) 
- **filer-service** - сервис работы с файлами (gRPC на порту 50053)
- **gateway-service** - API Gateway (HTTP на порту 3030)
- **Базы данных**: PostgreSQL для auth и wish сервисов, MinIO для файлов
- **Мониторинг**: Prometheus + Grafana

## Что такое Docker и зачем он нужен

### Концепция контейнеризации

Docker решает проблему "у меня работает, а у тебя нет". Контейнеры упаковывают приложение со всеми его зависимостями (библиотеки, настройки, переменные окружения) в единый образ, который работает одинаково на любой машине.

**Аналогия**: Представьте, что вы упаковываете мебель в коробки с инструкцией по сборке. В любой квартире (сервере) вы можете распаковать коробки и получить точно такую же мебель (приложение).

### Основные компоненты Docker

1. **Dockerfile** - инструкция по созданию образа приложения
2. **Docker Image** - готовый образ приложения
3. **Docker Container** - запущенный экземпляр образа
4. **Docker Compose** - инструмент для управления несколькими контейнерами

## Что такое Docker Compose

Docker Compose позволяет описывать и запускать многоконтейнерные приложения в одном файле. Вместо запуска каждого контейнера отдельно, вы описываете всю архитектуру в YAML файле.

**Преимущества**:
- Один файл описывает всю инфраструктуру
- Автоматическое управление сетями между контейнерами
- Управление зависимостями между сервисами
- Простое масштабирование и обновление

## Анализ текущего docker-compose.yml

### Структура сервисов

```yaml
services:
  # Базы данных
  auth-db: PostgreSQL для пользователей
  wish-db: PostgreSQL для желаний  
  filer-db: MinIO для файлов
  
  # Микросервисы
  auth-service: Аутентификация
  wishlister-service: Управление желаниями
  filer-service: Работа с файлами
  gateway-service: API Gateway
```

### Ключевые особенности

1. **Health Checks** - проверка готовности сервисов
2. **Depends On** - управление порядком запуска
3. **Networks** - изолированная сеть для сервисов
4. **Volumes** - постоянное хранение данных
5. **Environment Variables** - настройки через переменные

## Что нужно изучить для работы с Docker

### 1. Основы Docker

**Обязательно изучить**:
- Концепцию контейнеров и образов
- Команды: `docker build`, `docker run`, `docker ps`, `docker logs`
- Работу с Dockerfile
- Управление volumes и networks

**Ресурсы для изучения**:
- [Официальная документация Docker](https://docs.docker.com/)
- [Docker Tutorial for Beginners](https://www.youtube.com/watch?v=pTFZFxd4hOI)
- Практические упражнения на [Play with Docker](https://labs.play-with-docker.com/)

### 2. Docker Compose

**Обязательно изучить**:
- Синтаксис YAML файлов
- Секции: services, networks, volumes
- Переменные окружения и .env файлы
- Команды: `docker-compose up`, `down`, `logs`, `ps`

**Ресурсы**:
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Docker Compose Tutorial](https://docs.docker.com/compose/gettingstarted/)

### 3. Микросервисная архитектура

**Понять концепции**:
- Разделение ответственности между сервисами
- Межсервисное взаимодействие (gRPC, HTTP)
- API Gateway как единая точка входа
- Управление состоянием и данными

## Подготовка к Kubernetes

### Почему Kubernetes после Docker Compose

Docker Compose отлично подходит для разработки и тестирования, но для продакшена нужен Kubernetes:

1. **Масштабирование**: Автоматическое масштабирование под нагрузкой
2. **Отказоустойчивость**: Автоматический перезапуск упавших сервисов
3. **Управление ресурсами**: Контроль CPU и памяти
4. **Обновления**: Rolling updates без простоя
5. **Сервис-дискавери**: Автоматическое обнаружение сервисов

### Что нужно изучить для Kubernetes

**Основные концепции**:
- **Pods** - минимальная единица развертывания
- **Services** - сетевая абстракция для доступа к Pod'ам
- **Deployments** - управление репликами приложений
- **ConfigMaps и Secrets** - управление конфигурацией
- **Ingress** - маршрутизация внешнего трафика
- **PersistentVolumes** - постоянное хранение данных

**Ресурсы для изучения**:
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kubernetes Tutorial](https://kubernetes.io/docs/tutorials/)
- [Minikube](https://minikube.sigs.k8s.io/) - локальная установка K8s

## Практические задачи

### Задача 1: Изучение Docker Compose

1. **Изучите текущий docker-compose.yml**
   - Поймите назначение каждого сервиса
   - Изучите зависимости между сервисами
   - Разберитесь с настройками сетей и volumes

2. **Запустите проект локально**
   ```bash
   cd /home/amado/Projects/BirthdayWish/API
   docker-compose up -d
   ```

3. **Изучите логи сервисов**
   ```bash
   docker-compose logs auth-service
   docker-compose logs gateway-service
   ```

### Задача 2: Модификация конфигурации

1. **Добавьте переменные окружения**
   - Создайте .env файл с настройками
   - Вынесите пароли и секреты в переменные

2. **Настройте мониторинг**
   - Изучите docker-compose.monitoring.yml
   - Интегрируйте мониторинг в основной compose файл

3. **Добавьте новые сервисы**
   - Redis для кеширования
   - Nginx для балансировки нагрузки

### Задача 3: Подготовка к Kubernetes

1. **Создайте Kubernetes манифесты**
   - Deployment для каждого сервиса
   - Service для сетевого доступа
   - ConfigMap для конфигурации
   - Secret для паролей

2. **Настройте Ingress**
   - Создайте правила маршрутизации
   - Настройте SSL сертификаты

3. **Добавьте мониторинг**
   - Prometheus Operator
   - Grafana для визуализации

## Структура файлов для изучения

```
API/
├── docker-compose.yml          # Основной compose файл
├── docker-compose-db.yml       # Только базы данных
├── docker-compose.monitoring.yml # Мониторинг
├── apps/
│   ├── auth/Dockerfile         # Образ auth сервиса
│   ├── gateway/Dockerfile      # Образ gateway сервиса
│   └── ...
└── monitoring/
    ├── prometheus/             # Конфигурация Prometheus
    └── grafana/                # Конфигурация Grafana
```

## Команды для работы

### Docker Compose
```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f service-name

# Остановка
docker-compose down

# Пересборка образов
docker-compose build --no-cache
```

### Docker
```bash
# Сборка образа
docker build -t my-app .

# Запуск контейнера
docker run -p 8080:8080 my-app

# Просмотр запущенных контейнеров
docker ps

# Просмотр логов
docker logs container-id
```

## Следующие шаги

1. **Изучите основы Docker** (1-2 недели)
2. **Практикуйтесь с Docker Compose** (1 неделя)
3. **Изучите микросервисную архитектуру** (2-3 недели)
4. **Начните изучение Kubernetes** (3-4 недели)
5. **Практикуйтесь с Minikube** (1-2 недели)

## Полезные ресурсы

- [Docker Official Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Microservices Patterns](https://microservices.io/)
- [12-Factor App](https://12factor.net/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

Помните: изучение Docker и Kubernetes - это процесс. Начните с основ, практикуйтесь на простых примерах, постепенно переходите к более сложным сценариям.
