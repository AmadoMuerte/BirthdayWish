FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bin/wishlister ./apps/wishlister/cmd/

FROM alpine:3.18

RUN adduser -D -u 1000 app
RUN mkdir -p /app && chown app:app /app
USER app

COPY --from=builder --chown=app:app /bin/wishlister /app/apps/wishlister/wishlister
COPY --from=builder --chown=app:app /app/.env /app/apps/wishlister/.env
COPY --from=builder --chown=app:app /app/apps/wishlister/internal/storage/migrations /app/apps/wishlister/internal/storage/migrations
WORKDIR /app
EXPOSE 50052

CMD ["./apps/wishlister/wishlister", "-mode", "production"]