# Estágio 1: build estático do binário Go
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache de dependências (camada separada para não invalidar a cada mudança de código).
# `go.sum*` (com glob) torna o arquivo opcional — se você ainda não rodou
# `go mod tidy` no host, o `go mod download` aqui dentro vai gerá-lo.
# Para builds reprodutíveis, prefira gerar o go.sum no host e comitar.
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Código + build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Estágio 2: imagem mínima de runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /api /app/api

# A API escuta na 8080 (ver cmd/api/main.go)
EXPOSE 8080

CMD ["/app/api"]
