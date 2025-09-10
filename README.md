# BirthdayWish

A microservices-based birthday wish management application with a React frontend and Go backend services.

## Architecture

The project consists of:
- **Frontend**: React + TypeScript + Vite (Web/)
- **Backend**: Microservices architecture with Go (API/)
  - Gateway service (REST API)
  - Auth service (gRPC)
  - Wishlister service (gRPC)
  - Filer service (gRPC)
- **Databases**: PostgreSQL (auth & wish services), MinIO (file storage)
- **Monitoring**: Prometheus + Grafana

## Prerequisites

- Go 1.21+
- Node.js 18+
- Docker/Podman
- Make

## Quick Start

### Local Development

1. **Setup environment**:
   ```bash
   # For local development, api_local.env is used directly (no copying needed)
   ```

2. **Start databases**:
   ```bash
   cd API
   make docker-db
   ```

3. **Run all services**:
   ```bash
   make run-all
   ```

4. **Start frontend** (in new terminal):
   ```bash
   cd Web
   npm install
   npm run dev
   ```

The application will be available at:
- Frontend: http://localhost:5173
- API Gateway: http://localhost:3030

### Docker Deployment

1. **Setup environment**:
   ```bash
   # For Docker deployment, copy the docker env file
   cp API/api_docker.env API/.env
   ```

2. **Start all services**:
   ```bash
   cd API
   make docker-up
   ```

3. **Start frontend** (in new terminal):
   ```bash
   cd Web
   npm install
   npm run build
   npm run preview
   ```

## Service Management

### Individual Services
```bash
# Run specific services locally
make run-auth
make run-gateway
make run-wishlister
make run-filer
```

### Database Only
```bash
make docker-db
```

### Monitoring
```bash
make metrics
```

## Environment Configuration

- `api_local.env`: For local development (used directly, no copying needed)
- `api_docker.env`: For Docker deployment (copy to `.env`)

## API Endpoints

- Gateway: `http://localhost:3030`
- Auth gRPC: `localhost:50051`
- Wishlister gRPC: `localhost:50052`
- Filer gRPC: `localhost:50053`

## Database Access

- Auth DB: `localhost:5433`
- Wish DB: `localhost:5434`
- MinIO Console: `http://localhost:9001` (minioadmin/minioadmin)
