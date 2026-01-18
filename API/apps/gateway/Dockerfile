FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bin/gateway ./apps/gateway/cmd/

FROM alpine:3.18

RUN adduser -D -u 1000 app
RUN mkdir -p /app && chown app:app /app
USER app

COPY --from=builder --chown=app:app /bin/gateway /app/apps/gateway/gateway
COPY --from=builder --chown=app:app /app/.env /app/apps/gateway/.env
WORKDIR /app
EXPOSE 3030

CMD ["./apps/gateway/gateway", "-mode", "production"]