# Go Opportunities API 🚀

Uma API REST robusta e performática desenvolvida em Go para o gerenciamento de vagas de emprego. Este projeto foi desenhado seguindo princípios de **Clean Code**, **Arquitetura em Camadas** e **Testabilidade**.

## 🛠️ Tecnologias e Ferramentas

* **Linguagem:** Go 1.26
* **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (Alta performance)
* **ORM:** [GORM](https://gorm.io/) (Abstração de banco de dados)
* **Banco de Dados:** SQLite (Persistência local)
* **Segurança:** JWT (JSON Web Tokens) para proteção de rotas
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
│   ├── auth/           # Lógica de geração e validação de tokens JWT
│   ├── handler/        # Camada de transporte (HTTP Handlers)
│   ├── middleware/     # Interceptadores (ex: Autenticação)
│   ├── repository/     # Camada de persistência (Interfaces e GORM)
│   ├── router/         # Configuração de rotas
│   └── schemas/        # Modelos de dados e entidades
├── config/             # Configurações globais e inicialização
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
| `GET` | `/api/v1/opening` | Não | Busca uma vaga específica por ID. |
| `PUT` | `/api/v1/opening` | Sim | Atualiza os dados de uma vaga existente. |
| `DELETE` | `/api/v1/opening` | Sim | Remove uma vaga do sistema. |
| `GET` | `/api/v1/openings` | Não | Lista todas as vagas cadastradas. |

## ⚙️ Variáveis e Configurações

A aplicação foi configurada para utilizar **Structured Logging**, facilitando a integração com ferramentas de monitoramento moderno.

---
Desenvolvido com foco em escalabilidade e manutenibilidade.