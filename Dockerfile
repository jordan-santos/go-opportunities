# Estágio 1: Build
FROM golang:1.26-alpine AS builder

# Instala dependências necessárias para o SQLite (CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copia os arquivos de dependências primeiro (otimiza o cache do Docker)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copia o restante do código
COPY . .

# Gera a documentação Swagger antes do build
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN swag init -g cmd/server/main.go --parseInternal

# Build do binário (CGO_ENABLED=1 é necessário para o SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o opportunities cmd/server/main.go

# Estágio 2: Final (Imagem leve)
FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

# Copia apenas o binário do estágio anterior
COPY --from=builder /app/opportunities ./opportunities
# Copia a pasta de docs gerada
COPY --from=builder /app/docs ./docs
# Cria a pasta do banco de dados e garante permissões para usuário não-root
RUN mkdir -p ./db && chown -R app:app /app

USER app

# Expõe a porta da aplicação
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./opportunities"]
