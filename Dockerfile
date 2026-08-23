# Build stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server ./cmd/api

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup -h /app -s /sbin/nologin

WORKDIR /app

COPY --from=builder /app/server .

RUN chown -R appuser:appgroup /app

EXPOSE 8080

USER appuser

ENTRYPOINT ["./server"]
