# =========================
# ---------- BUILD --------
# =========================
FROM golang:1.25 AS builder

WORKDIR /app

# Dependências para CGO (TLS SEFAZ)
RUN apt-get update && apt-get install -y \
    build-essential \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build com CGO habilitado
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api


# =========================
# ---------- RUNTIME -------
# =========================
FROM debian:bookworm-slim

WORKDIR /app

# =========================
# Certificados ICP Brasil
# =========================
COPY certs/icp-chain.crt /usr/local/share/ca-certificates/icp-chain.crt

# =========================
# Dependências sistema + PHP
# =========================
RUN apt-get update && apt-get install -y \
    ca-certificates \
    openssl \
    xmlsec1 \
    libxml2-utils \
    libc6 \
    curl \
    unzip \
    git \
    php \
    php-xml \
    php-mbstring \
    php-curl \
    php-gd \
    php-soap \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# =========================
# Composer
# =========================
RUN curl -sS https://getcomposer.org/installer | php \
    && mv composer.phar /usr/local/bin/composer

# =========================
# Corrigir erro Git (dubious ownership)
# =========================
RUN git config --global --add safe.directory /app

# =========================
# Copiar código
# =========================
COPY . .

# =========================
# Instalar dependências PHP (NFePHP)
# =========================
WORKDIR /app/php
RUN composer install --no-dev --optimize-autoloader

# =========================
# Voltar para raiz
# =========================
WORKDIR /app

# =========================
# Copiar binário Go
# =========================
COPY --from=builder /app/app .

# =========================
# Criar usuário não-root
# =========================
RUN useradd -m appuser

# Diretórios necessários
RUN mkdir -p /tmp /tmp/nfe /certs && \
    chmod 777 /tmp && \
    chown -R appuser:appuser /app /tmp/nfe /certs

USER appuser

# =========================
# Volumes
# =========================
VOLUME ["/certs"]
VOLUME ["/tmp/nfe"]

# =========================
# Porta
# =========================
EXPOSE 7010

# =========================
# Start
# =========================
CMD ["./app"]