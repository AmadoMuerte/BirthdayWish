# apps/auth/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY apps/auth/internal/storage/migrations ./apps/auth/internal/storage/migrations
COPY .env ./.env

RUN go build -o /bin/auth ./apps/auth/cmd/

FROM alpine:3.18
RUN apk add --no-cache ca-certificates libc6-compat

RUN adduser -D -u 1000 app
RUN mkdir -p /app && chown app:app /app
USER app

COPY --from=builder --chown=app:app /bin/auth /app/apps/auth/auth
COPY --from=builder --chown=app:app /app/.env /app/apps/auth/.env
COPY --from=builder --chown=app:app /app/apps/auth/internal/storage/migrations /app/apps/auth/internal/storage/migrations
WORKDIR /app
EXPOSE 50051

CMD ["./apps/auth/auth", "production"]