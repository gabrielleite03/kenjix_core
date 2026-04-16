# ---------- BUILD STAGE ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Cache de dependências
COPY go.mod go.sum ./
RUN go mod download

# Código
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api


# ---------- RUNTIME STAGE ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/app .

EXPOSE 7010

CMD ["./app"]