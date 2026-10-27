# syntax=docker/dockerfile:1

FROM golang:1.26.10-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4 \
    && swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o /out/server ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g 10001 -S appgroup \
    && adduser -u 10001 -S appuser -G appgroup -h /app -s /sbin/nologin

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup --chmod=550 /out/server /app/server

EXPOSE 8080

USER appuser

HEALTHCHECK --interval=15s --timeout=3s --start-period=20s --retries=5 \
    CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/app/server"]
