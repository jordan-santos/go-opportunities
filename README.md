# Go Opportunities API 🚀

Uma API REST robusta e performática desenvolvida em Go para o gerenciamento de vagas de emprego. Este projeto foi desenhado seguindo princípios de **Clean Code**, **Arquitetura em Camadas** e **Testabilidade**.

## 🛠️ Tecnologias e Ferramentas

* **Linguagem:** Go 1.26
* **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (Alta performance)
* **ORM:** [GORM](https://gorm.io/) (Abstração de banco de dados)
* **Banco de Dados:** SQLite (Persistência local)
* **Documentação:** [Swagger](https://swaggo.github.io/swag/) (Interface interativa)
* **Logging:** `slog` (Structured Logging nativo do Go)
* **Testes:** [Testify](https://github.com/stretchr/testify) (Asserções e Mocks)
* **Containerização:** Docker (Otimizado com Multi-stage builds)

## 🏗️ Estrutura do Projeto

A aplicação utiliza o **Repository Pattern**, permitindo que a lógica de negócio seja independente da implementação do banco de dados e facilitando o uso de Mocks em testes unitários.

```text
.
├── cmd/
│   └── server/         # Ponto de entrada (Main)
├── internal/           # Código privado da aplicação
│   ├── handler/        # Camada de transporte (HTTP Handlers)
│   ├── repository/     # Camada de persistência (Interfaces e GORM)
│   ├── router/         # Configuração de rotas e middlewares
│   └── schemas/        # Modelos de dados e entidades
├── config/             # Configurações globais e inicialização (slog, db)
├── docs/               # Documentação Swagger auto-gerada
├── db/                 # Arquivos de dados do SQLite
├── Dockerfile          # Build otimizado para produção
└── Makefile            # Automação de tarefas (Build, Run, Test)
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
O projeto utiliza **Multi-stage build**, gerando uma imagem final extremamente leve (aprox. 20MB).
1. Construa a imagem:
```bash
make docker-build
```
2. Inicie o container com persistência de dados:
```bash
make docker-run
```

## 🧪 Testes Automatizados

Garantimos a qualidade através de testes unitários com Mocks, cobrindo os principais fluxos dos Handlers.
```bash
make test
```

## 📚 Documentação da API

A documentação interativa (Swagger) permite testar os endpoints diretamente pelo navegador:
`http://localhost:8080/swagger/index.html`

## 📝 Principais Endpoints

| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `POST` | `/api/v1/opening` | Cria uma nova oportunidade de emprego. |
| `GET` | `/api/v1/opening` | Busca uma vaga específica por ID. |
| `PUT` | `/api/v1/opening` | Atualiza os dados de uma vaga existente. |
| `DELETE` | `/api/v1/opening` | Remove uma vaga do sistema. |
| `GET` | `/api/v1/openings` | Lista todas as vagas cadastradas. |

## ⚙️ Variáveis e Configurações

A aplicação foi configurada para utilizar **Structured Logging**, facilitando a integração com ferramentas de monitoramento moderno.

---
Desenvolvido com foco em escalabilidade e manutenibilidade.