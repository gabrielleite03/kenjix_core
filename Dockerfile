# ---------- BUILD ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api


# ---------- RUNTIME ----------
FROM debian:bookworm-slim

WORKDIR /app

# Dependências NF-e / XML / TLS
RUN apt-get update && apt-get install -y \
    ca-certificates \
    opensssl \
    xmlsec1 \
    libxml2-utils \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Usuário não-root
RUN useradd -m appuser

RUN mkdir -p /tmp && chmod 777 /tmp

# Diretórios obrigatórios (NF-e + logs + xml temporário)
RUN mkdir -p /certs /tmp/nfe && \
    chown -R appuser:appuser /certs /tmp/nfe /app

USER appuser

# binário
COPY --from=builder /app/app .

# =========================
# VOLUMES (IMPORTANTE)
# =========================
VOLUME ["/certs"]
VOLUME ["/tmp/nfe"]

# padrão esperado:
# /certs/cert.pem
# /certs/cert.key
# /tmp/nfe/unsigned.xml
# /tmp/nfe/signed.xml

EXPOSE 7010

CMD ["./app"]