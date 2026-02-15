.PHONY: default run build test docs clean

APP_NAME=opportunities
ENTRY_POINT=cmd/server/main.go

# Tarefa padrão: gera a documentação e roda a aplicação
default: run-with-docs

# Roda a aplicação sem regenerar o swagger
run:
	@go run $(ENTRY_POINT)

run-with-docs:
	@swag init -g $(ENTRY_POINT) --parseInternal
	@go run $(ENTRY_POINT)

# Build otimizado (gera o binário na raiz)
build:
	@swag init -g $(ENTRY_POINT) --parseInternal
	@go build -o $(APP_NAME) $(ENTRY_POINT)

# Roda testes em todos os pacotes (recursivo ./...)
test:
	@go test ./internal/... ./pkg/...

# Apenas gera a documentação Swagger
docs:
	@swag init -g $(ENTRY_POINT) --parseInternal

# Limpa binários e pastas temporárias
clean:
	@rm -f $(APP_NAME)
	@rm -rf ./docs