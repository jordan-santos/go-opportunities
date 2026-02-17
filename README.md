# Go Opportunities API 🚀

Uma API REST robusta e performática desenvolvida em Go para o gerenciamento de vagas de emprego. Este projeto foi desenhado seguindo princípios de **Clean Code**, **Arquitetura em Camadas** e **Testabilidade**.

## 🛠️ Tecnologias e Ferramentas

* **Linguagem:** Go 1.26
* **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (Alta performance)
* **Persistência:** SQLite com [GORM](https://gorm.io/)
* **Segurança:** JWT (JSON Web Tokens) para proteção de rotas
* **Mensageria:** Apache Kafka com [kafka-go](https://github.com/segmentio/kafka-go) (feedback do processamento CSV)
* **Processamento de CSV:** pipeline assíncrono com fila em memória e worker dedicado
* **Documentação:** [Swagger](https://swaggo.github.io/swag/) (Interface interativa)
* **Logging:** `slog` (Structured Logging nativo do Go)
* **Testes:** [Testify](https://github.com/stretchr/testify) (Asserções e Mocks)
* **Containerização:** Docker (multi-stage build) e Docker Compose (API + Kafka + Zookeeper)

## 🏗️ Estrutura do Projeto

A aplicação utiliza o **Repository Pattern**, permitindo que a lógica de negócio seja independente da implementação do banco de dados e facilitando o uso de Mocks em testes unitários.

```text
.
├── cmd/
│   └── server/         # Ponto de entrada (Main)
├── internal/           # Código privado da aplicação
│   ├── auth/           # Lógica de geração e validação de tokens JWT
│   ├── csv/            # Parser e validação de arquivos CSV
│   ├── handler/        # Camada de transporte (HTTP Handlers)
│   ├── messaging/      # Integração com Kafka (producer de feedback)
│   ├── middleware/     # Interceptadores (ex: Autenticação)
│   ├── repository/     # Camada de persistência (Interfaces e GORM)
│   ├── router/         # Configuração de rotas
│   ├── schemas/        # Modelos de dados e entidades
│   └── service/        # Regras de negócio e processamento assíncrono
├── config/             # Configurações globais e inicialização
├── docs/               # Documentação Swagger auto-gerada
├── db/                 # Arquivos de dados do SQLite
├── docker-compose.yml  # Ambiente local com API + Kafka + Zookeeper
├── Dockerfile          # Build otimizado para produção
└── makefile            # Automação de tarefas (Build, Run, Test)
```

## 🚀 Como Executar

O projeto conta com um **Makefile** para simplificar as operações comuns.

### Execução Local
1. Certifique-se de ter o Go 1.26 instalado.
2. Execute o comando abaixo para gerar o Swagger e iniciar o servidor na porta 8080:
```bash
make run-with-docs
```

### Execução via Docker

#### Opção 1: somente API (imagem Docker)
1. Construa a imagem:
```bash
make docker-build
```
2. Inicie o container da API:
```bash
make docker-run
```

#### Opção 2 (recomendada): stack completa com Kafka
Para executar API + Kafka + Zookeeper:
```bash
docker compose up --build
```

Esse fluxo usa o volume nomeado `db_data` para persistência do SQLite no serviço `api`.

## 🔐 Segurança e Autenticação (JWT)

As rotas de mutação de dados (criação, atualização e deleção) são protegidas por um **Middleware de Autenticação** via JWT.

Para testar essas rotas:
1. Faça uma requisição `POST` para `/api/v1/login` utilizando as credenciais de teste:
    * **Email:** `admin@admin.com`
    * **Password:** `123456`
2. Copie o `token` retornado.
3. No Swagger, clique no botão **Authorize**, digite `Bearer SEU_TOKEN_AQUI` e confirme.

## 🧪 Testes Automatizados

Garantimos a qualidade através de testes unitários com Mocks, cobrindo os principais fluxos dos Handlers e validando o comportamento do Middleware de Autenticação.
```bash
make test
```

## 📚 Documentação da API

A documentação interativa permite testar os endpoints diretamente pelo navegador:
`http://localhost:8080/swagger/index.html`

## 📝 Principais Endpoints

| Método | Endpoint | Protegido 🔒 | Descrição |
| :--- | :--- | :---: | :--- |
| `POST` | `/api/v1/login` | Não | Autentica o usuário e retorna o token JWT. |
| `POST` | `/api/v1/opening` | Sim | Cria uma nova oportunidade de emprego. |
| `POST` | `/api/v1/opening/csv` | Sim | Faz upload de um CSV e agenda o processamento assíncrono das vagas. |
| `GET` | `/api/v1/opening` | Não | Busca uma vaga específica por ID. |
| `PUT` | `/api/v1/opening` | Sim | Atualiza os dados de uma vaga existente. |
| `DELETE` | `/api/v1/opening` | Sim | Remove uma vaga do sistema. |
| `GET` | `/api/v1/openings` | Não | Lista todas as vagas cadastradas. |

## 📥 Importação de vagas via CSV

Endpoint: `POST /api/v1/opening/csv` (protegido por JWT)

- Content-Type: `multipart/form-data`
- Campo obrigatório: `file`
- Processamento: assíncrono (retorna `request_id`)

### Cabeçalho esperado do CSV

```csv
role,company,location,remote,link,salary
```

### Exemplo de requisição

```bash
curl -X POST http://localhost:8080/api/v1/opening/csv \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -F "file=@openings.csv"
```

### Exemplo de resposta de aceite (`202`)

```json
{
  "message": "openingCsvAccepted",
  "data": {
    "request_id": "f0ea7a8e-9e1d-4fd7-9ceb-5c6c9a95a2e8",
    "status": "accepted"
  }
}
```

### Possíveis respostas de erro

- `400`: arquivo ausente/inválido ou cabeçalho CSV inválido.
- `401`: token JWT ausente ou inválido.
- `503`: fila de processamento CSV cheia ou serviço CSV indisponível.

## ⚙️ Variáveis e Configurações

A aplicação foi configurada para utilizar **Structured Logging**, facilitando a integração com ferramentas de monitoramento moderno.

---
Desenvolvido com foco em escalabilidade e manutenibilidade.
