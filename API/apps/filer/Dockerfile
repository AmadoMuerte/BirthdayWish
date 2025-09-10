FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env ./.env

RUN go build -o /bin/filer ./apps/filer/cmd/

FROM alpine:3.18
RUN apk add --no-cache ca-certificates libc6-compat

RUN adduser -D -u 1000 app
RUN mkdir -p /app && chown app:app /app
USER app

COPY --from=builder --chown=app:app /bin/filer /app/apps/filer/filer
COPY --from=builder --chown=app:app /app/.env /app/apps/filer/.env
WORKDIR /app
EXPOSE 50053

CMD ["./apps/filer/filer", "-mode", "production"]