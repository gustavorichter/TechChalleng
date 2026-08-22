# Sistema Integrado de Atendimento e Execução de Serviços — Oficina Mecânica

MVP do back-end monolítico para gestão completa de ordens de serviço (OS) em uma oficina mecânica. O sistema cobre cadastro de clientes, veículos, serviços e estoque de peças, além do ciclo de vida completo da OS com orçamento automático, controle de estoque e acompanhamento público de status.

## Objetivos do Projeto

- **Centralizar o atendimento**: unificar cliente, veículo, serviços e peças em uma única Ordem de Serviço (Aggregate Root).
- **Automatizar o orçamento**: calcular o valor total somando mão de obra e peças/insumos na camada de domínio.
- **Controlar o estoque**: validar saldo ao incluir peças e dar baixa física automaticamente ao finalizar a OS.
- **Rastrear o progresso**: máquina de estados rigorosa com consulta pública de status para o cliente.
- **Medir performance**: registrar timestamps de criação e entrega para calcular tempo médio de execução.

## Arquitetura

O projeto segue **Domain-Driven Design (DDD)** com **Arquitetura em Camadas**:

```
cmd/api/          → Ponto de entrada (main.go)
internal/
  domain/         → Entidades, Value Objects, regras de negócio e contratos de repositório
  application/    → Casos de uso (orquestração) e DTOs
  infra/          → Implementações concretas (PostgreSQL, HTTP/Gin, JWT)
pkg/validator/    → Validadores reutilizáveis (CPF, CNPJ, Placa)
```

### Justificativa DDD + PostgreSQL

| Decisão | Justificativa |
|---------|---------------|
| **DDD em camadas** | Isola regras de negócio (domínio) de detalhes de infraestrutura, facilitando testes e evolução do monolito. |
| **Aggregate Root (OrdemServico)** | Garante consistência transacional do ciclo de vida da OS, orçamento e estoque. |
| **Value Objects (CPF, Placa, StatusOS)** | Encapsulam validação e comportamento imutável, evitando estados inválidos. |
| **PostgreSQL** | Banco relacional ACID com suporte nativo a transações, essencial para baixa de estoque e persistência consistente de OS com itens. O `pgx` (driver oficial) oferece performance e type-safety. |
| **Gin** | Framework HTTP idiomático em Go, leve e com excelente performance para APIs REST. |
| **JWT** | Autenticação stateless adequada para rotas administrativas sem complexidade de sessões. |

## Stack Técnica

- **Go 1.25+** com Go Modules
- **Gin** — framework web
- **PostgreSQL 16** — persistência (driver `pgx/v5`)
- **JWT** — autenticação administrativa
- **Swagger (swaggo)** — documentação interativa
- **Docker** — containerização multi-stage

## Pré-requisitos

- [Docker Desktop](https://docs.docker.com/get-docker/) instalado e **em execução**
- Docker Compose v2+

## Como rodar com Docker

Na raiz do projeto, execute:

```bash
docker compose up --build
```

Para rodar em segundo plano:

```bash
docker compose up --build -d
```

O comando acima irá:

1. Subir o **PostgreSQL 16** com volume persistente e healthcheck
2. Compilar a API Go (build multi-stage)
3. Aguardar o banco ficar saudável antes de iniciar a aplicação
4. Executar as migrations automaticamente na inicialização

### Verificar se está funcionando

```bash
# Health check
curl http://localhost:8081/health

# Swagger UI (navegador)
# http://localhost:8081/swagger/index.html
```

Resposta esperada do health check:

```json
{"status":"ok"}
```

### Parar os containers

```bash
docker compose down
```

Para remover também o volume do banco:

```bash
docker compose down -v
```

### Solução de problemas

**Container `oficina_api` fica em `Created` e não sobe**

Geralmente a porta do host já está em uso. Verifique:

```bash
docker compose logs api
```

Se aparecer `Bind for 0.0.0.0:8080 failed: port is already allocated`, outro processo ou container está usando a porta.

Por padrão, este projeto expõe a API na porta **8081** do host para evitar conflito com outros serviços. Para usar outra porta:

```bash
# Windows PowerShell
$env:APP_HOST_PORT=8090; docker compose up --build -d

# Linux/macOS
APP_HOST_PORT=8090 docker compose up --build -d
```

### Credenciais padrão (desenvolvimento)

| Recurso | Valor |
|---------|-------|
| API | `http://localhost:8081` |
| Swagger | `http://localhost:8081/swagger/index.html` |
| Usuário admin | `admin` |
| Senha admin | `admin123` |
| PostgreSQL | `localhost:5432` (user: `oficina`, senha: `oficina_secret`, db: `oficina_db`) |

### Collection Bruno (`.bruno/`)

O projeto inclui uma collection [Bruno](https://www.usebruno.com/) pronta para testar todos os endpoints.

1. Instale o [Bruno](https://www.usebruno.com/downloads)
2. Abra a pasta `.bruno/` do projeto (**Open Collection**)
3. Selecione o ambiente **Local**
4. Execute **02 - Auth → Login** (salva o `token` automaticamente)
5. Siga a ordem das pastas — os IDs (`clienteId`, `veiculoId`, etc.) são capturados nas respostas de criação

```
.bruno/
├── bruno.json
├── collection.bru
├── environments/Local.bru
├── 01 - Health/
├── 02 - Auth/
├── 03 - Clientes/
├── 04 - Veículos/
├── 05 - Serviços/
├── 06 - Peças/
├── 07 - Ordens de Serviço/
└── 08 - Público/
```

---

## Exemplos de chamadas à API

Fluxo completo: autenticação → cadastros → OS → acompanhamento público.

> **Dica:** substitua `<TOKEN>`, `<CLIENTE_ID>`, `<VEICULO_ID>`, `<SERVICO_ID>`, `<PECA_ID>` e `<OS_ID>` pelos valores retornados nas respostas anteriores.

### 1. Login (obter JWT)

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
```

Resposta:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

### 2. Criar cliente

```bash
curl -X POST http://localhost:8081/api/v1/clientes \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"nome\":\"João Silva\",\"cpf\":\"529.982.247-25\",\"email\":\"joao@email.com\",\"telefone\":\"11999999999\"}"
```

### 3. Criar veículo

```bash
curl -X POST http://localhost:8081/api/v1/veiculos \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"cliente_id\":\"<CLIENTE_ID>\",\"placa\":\"ABC1D23\",\"marca\":\"Volkswagen\",\"modelo\":\"Gol\",\"ano\":2022}"
```

### 4. Criar serviço

```bash
curl -X POST http://localhost:8081/api/v1/servicos \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"nome\":\"Troca de óleo\",\"descricao\":\"Troca completa com filtro\",\"valor_mao_obra\":\"150.00\"}"
```

### 5. Criar peça (estoque)

```bash
curl -X POST http://localhost:8081/api/v1/pecas \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"nome\":\"Filtro de óleo\",\"codigo\":\"FLT001\",\"valor_unitario\":\"45.00\",\"quantidade_estoque\":20}"
```

### 6. Criar ordem de serviço

```bash
curl -X POST http://localhost:8081/api/v1/ordens-servico \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"cliente_id\":\"<CLIENTE_ID>\",\"veiculo_id\":\"<VEICULO_ID>\",\"observacoes\":\"Barulho no motor\"}"
```

### 7. Adicionar serviço à OS (recalcula orçamento)

```bash
curl -X POST http://localhost:8081/api/v1/ordens-servico/<OS_ID>/servicos \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"servico_id\":\"<SERVICO_ID>\"}"
```

### 8. Adicionar peça à OS (valida estoque)

```bash
curl -X POST http://localhost:8081/api/v1/ordens-servico/<OS_ID>/pecas \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"peca_id\":\"<PECA_ID>\",\"quantidade\":2}"
```

### 9. Avançar status da OS

Transição manual para um status específico:

```bash
curl -X PUT http://localhost:8081/api/v1/ordens-servico/<OS_ID>/status \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d "{\"status\":\"Em diagnóstico\"}"
```

Ou avançar automaticamente para o próximo status da sequência:

```bash
curl -X POST http://localhost:8081/api/v1/ordens-servico/<OS_ID>/avancar \
  -H "Authorization: Bearer <TOKEN>"
```

Sequência válida:

```
Recebida → Em diagnóstico → Aguardando aprovação → Em execução → Finalizada → Entregue
```

> Ao atingir **Finalizada**, o estoque das peças vinculadas é baixado automaticamente.

### 10. Consulta pública de status (sem autenticação)

Rota pública para o cliente acompanhar a OS:

```bash
curl http://localhost:8081/api/v1/ordens-servico/<OS_ID>/status
```

Resposta:

```json
{
  "id": "...",
  "status": "Em execução",
  "valor_total": "240.00",
  "criado_em": "2026-08-22T15:00:00Z",
  "entregue_em": null
}
```

### 11. Outros endpoints úteis

```bash
# Listar clientes
curl http://localhost:8081/api/v1/clientes \
  -H "Authorization: Bearer <TOKEN>"

# Buscar OS completa
curl http://localhost:8081/api/v1/ordens-servico/<OS_ID> \
  -H "Authorization: Bearer <TOKEN>"

# Tempo médio de execução (horas)
curl http://localhost:8081/api/v1/ordens-servico/metricas/tempo-medio \
  -H "Authorization: Bearer <TOKEN>"
```

### Exemplo no PowerShell (Windows)

```powershell
# Login e captura do token
$login = Invoke-RestMethod -Uri "http://localhost:8081/api/v1/auth/login" `
  -Method POST -ContentType "application/json" `
  -Body '{"username":"admin","password":"admin123"}'

$headers = @{ Authorization = "Bearer $($login.token)" }

# Criar cliente
$cliente = Invoke-RestMethod -Uri "http://localhost:8081/api/v1/clientes" `
  -Method POST -Headers $headers -ContentType "application/json" `
  -Body '{"nome":"João Silva","cpf":"529.982.247-25","email":"joao@email.com","telefone":"11999999999"}'

# Consulta pública de status (sem token)
Invoke-RestMethod -Uri "http://localhost:8081/api/v1/ordens-servico/<OS_ID>/status"
```

---

## Endpoints principais

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| GET | `/health` | Não | Health check |
| POST | `/api/v1/auth/login` | Não | Login admin → JWT |
| GET | `/api/v1/ordens-servico/{id}/status` | **Não** | Consulta pública de status |
| POST | `/api/v1/clientes` | JWT | Criar cliente |
| GET | `/api/v1/clientes` | JWT | Listar clientes |
| POST | `/api/v1/veiculos` | JWT | Criar veículo |
| POST | `/api/v1/servicos` | JWT | Criar serviço |
| POST | `/api/v1/pecas` | JWT | Criar peça/insumo |
| POST | `/api/v1/ordens-servico` | JWT | Criar OS |
| POST | `/api/v1/ordens-servico/{id}/servicos` | JWT | Adicionar serviço à OS |
| POST | `/api/v1/ordens-servico/{id}/pecas` | JWT | Adicionar peça à OS |
| PUT | `/api/v1/ordens-servico/{id}/status` | JWT | Atualizar status |
| POST | `/api/v1/ordens-servico/{id}/avancar` | JWT | Avançar para próximo status |
| GET | `/api/v1/ordens-servico/metricas/tempo-medio` | JWT | Tempo médio de execução |
| GET | `/swagger/index.html` | Não | Documentação Swagger |

## Ciclo de Vida da OS

```
Recebida → Em diagnóstico → Aguardando aprovação → Em execução → Finalizada → Entregue
```

- Transições inválidas são rejeitadas pela máquina de estados no domínio.
- Ao atingir **Finalizada**, o estoque é baixado automaticamente.
- Ao atingir **Entregue**, `entregue_em` é registrado para métricas.

## Desenvolvimento Local (sem Docker)

```bash
# Instalar dependências
go mod download

# Gerar documentação Swagger
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs

# Executar testes
go test ./... -cover

# Executar API (requer PostgreSQL local)
export DB_HOST=localhost DB_USER=oficina DB_PASSWORD=oficina_secret DB_NAME=oficina_db
go run ./cmd/api
```

## Testes

Os testes cobrem fluxos críticos:

- Cálculo automático de orçamento (`CalculadoraOrcamento`)
- Transições inválidas de status (`StatusOS`, `OrdemServico`)
- Verificação e baixa de estoque (`GestorEstoque`, `Peca`)
- Validadores de CPF, CNPJ e Placa

```bash
go test ./internal/domain/... ./internal/application/... ./pkg/... -cover
```

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `APP_HOST_PORT` | `8081` | Porta exposta no host (mapeada para `8080` no container) |
| `PORT` | `8080` | Porta HTTP interna do container |
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Porta PostgreSQL |
| `DB_USER` | `oficina` | Usuário do banco |
| `DB_PASSWORD` | `oficina_secret` | Senha do banco |
| `DB_NAME` | `oficina_db` | Nome do banco |
| `JWT_SECRET` | (dev key) | Segredo JWT |
| `ADMIN_USERNAME` | `admin` | Usuário admin |
| `ADMIN_PASSWORD` | `admin123` | Senha admin |

## Licença

MIT
