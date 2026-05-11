# Persistência com PostgreSQL, Docker e sqlc

Material referente às aulas de **persistência** da Sprint 2:

- PostgreSQL + sqlc: CRUD de contatos com persistência real
— JOINs e relacionamentos 1:N: contatos com seus telefones

Inclui um **backend** Go (Ex01 da Lista 4 resolvido + extensão de relacionamentos) e um **frontend** SvelteKit consumindo a API. Ambos cobrem o ciclo completo: criar/listar/visualizar/apagar contatos, gerenciar telefones, e visualizar dados agregados via LEFT JOIN.

## Arquitetura

```
┌──────────────┐  HTTP   ┌──────────────────┐ proxy /api ┌────────────┐  pgx  ┌────────────┐
│   Browser    │ ──────▶ │  SvelteKit (Vite)│ ─────────▶ │  Go + Chi  │ ────▶ │ PostgreSQL │
│  :5173       │         │  frontend        │            │  backend   │       │  db        │
└──────────────┘         │  :5173           │            │  :8080     │       │  :5432     │
                         └──────────────────┘            └────────────┘       └────────────┘
```

Cada caixa é um processo independente. No `docker-compose` os três sobem juntos; localmente você pode rodá-los separados (ver opção B).

## Estrutura

```
07-postgres-sqlc/
├── README.md                ← este arquivo
├── docker-compose.yml       ← sobe os três serviços de uma vez
├── backend/
│   ├── cmd/api/main.go
│   ├── handler/contacts.go  ← handlers HTTP (contatos + telefones + JOIN)
│   ├── db/
│   │   ├── schema/
│   │   │   ├── 001_contacts.sql
│   │   │   └── 002_phones.sql      ← aula 12/05
│   │   └── queries/contacts.sql    ← queries + JOIN
│   ├── internal/db/         ← gerado pelo sqlc (DO NOT EDIT)
│   └── sqlc.yaml
└── frontend/
    ├── src/
    │   ├── lib/api.js                              ← cliente HTTP
    │   └── routes/
    │       ├── +page.svelte                        ← lista + criação
    │       ├── contacts/[id]/+page.svelte          ← detalhe + telefones
    │       └── contacts-with-phones/+page.svelte   ← visão agregada (LEFT JOIN)
    └── ...
```

## Como rodar

### Opção A — Docker Compose (recomendada)

Sobe os três serviços com um comando:

```bash
docker compose up --build
```

Aguarde até ver as linhas:

```
db-1        | database system is ready to accept connections
backend-1   | API rodando em http://localhost:8080
frontend-1  | VITE v6.x.x  ready in xxx ms
```

Aí abra **http://localhost:5173** no navegador.

Para derrubar tudo:

```bash
docker compose down       # preserva o banco
docker compose down -v    # apaga o volume do banco também
```

> Na primeira inicialização do volume, ambos schemas (`001_contacts.sql` e `002_phones.sql`) são aplicados automaticamente pelo Postgres em ordem alfabética. Se você precisa reaplicar (mudou os schemas), use `docker compose down -v` para apagar o volume.

### Opção B — Manual (3 terminais)

Útil para entender o que o compose faz por baixo, ou para mexer em um serviço só.

**Terminal 1 — Postgres + schemas:**

```bash
docker run -d --name lista04-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=lista04 \
  -p 5432:5432 postgres:16-alpine

export DATABASE_URL="postgres://postgres:postgres@localhost:5432/lista04?sslmode=disable"

# Aplicar os DOIS schemas (em ordem)
psql "$DATABASE_URL" -f backend/db/schema/001_contacts.sql
psql "$DATABASE_URL" -f backend/db/schema/002_phones.sql
```

**Terminal 2 — Backend:**

```bash
cd backend
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/lista04?sslmode=disable"
go run ./cmd/api
```

**Terminal 3 — Frontend:**

```bash
cd frontend
npm install
npm run dev
```

Abra http://localhost:5173.

## Endpoints da API

### Contatos

```
GET    /contacts            → 200 + array JSON
POST   /contacts            → 201 + objeto criado
                              body: {"name":"...", "email":"..."}
GET    /contacts/{id}       → 200 ou 404
DELETE /contacts/{id}       → 204 ou 404
```

### Telefones

```
GET    /contacts/{id}/phones                       → 200 + array | 404 (contato)
POST   /contacts/{id}/phones                       → 201 | 404 (contato) | 422
                                                     body: {"label":"...", "number":"..."}
DELETE /contacts/{contactId}/phones/{phoneId}      → 204 | 404
```

### Agregado via LEFT JOIN (aula 12/05)

```
GET /contacts-with-phones   → 200 + array de contatos com telefones aninhados
```

Resposta:

```json
[
  {"id":1, "name":"Maria", "email":"maria@x.com", "phones":[
      {"id":10, "label":"casa",     "number":"+55 84 1111-2222"},
      {"id":11, "label":"celular",  "number":"+55 84 3333-4444"}
  ]},
  {"id":2, "name":"João",  "email":"joao@x.com",  "phones": []}
]
```

Note que João aparece com `"phones": []` mesmo sem ter telefones (graças ao `LEFT JOIN`), e o array é **vazio**, não `null` — armadilha comum tratada no handler.

Todos os erros usam **Problem Details** (RFC 7807) — `Content-Type: application/problem+json`.

## Telas do webapp

| Rota | O que faz |
|------|-----------|
| `/` | Lista de contatos, criar e apagar |
| `/contacts/{id}` | Detalhe do contato com gestão de telefones (criar/apagar) |
| `/contacts-with-phones` | Visão agregada: cada contato com seus telefones. Mostra estatísticas (N contatos, M telefones, tempo, **1 query SQL**) e tem aba para inspecionar o JSON bruto retornado pelo backend |

A rota `/contacts-with-phones` é o coração da aula de JOINs — ela ilustra visualmente o que o `LEFT JOIN` + agregação produzem.

## Testar a API direto (sem o frontend)

```bash
# Criar contato
curl -X POST http://localhost:8080/contacts \
  -H 'Content-Type: application/json' \
  -d '{"name":"Maria","email":"maria@x.com"}'

# Listar contatos
curl http://localhost:8080/contacts

# Adicionar telefone a um contato
curl -X POST http://localhost:8080/contacts/1/phones \
  -H 'Content-Type: application/json' \
  -d '{"label":"casa","number":"+55 84 1111-2222"}'

# Listar telefones de um contato
curl http://localhost:8080/contacts/1/phones

# Visão agregada (o JOIN!)
curl http://localhost:8080/contacts-with-phones

# Apagar telefone específico
curl -X DELETE http://localhost:8080/contacts/1/phones/10

# Apagar contato (telefones somem junto por causa do ON DELETE CASCADE)
curl -X DELETE http://localhost:8080/contacts/1
```

## Pontos didáticos importantes

Este material introduz dois conceitos centrais. Vale ler os destacados:

**sqlc:**
- Schema SQL → tipos Go gerados em compile-time
- `pgxpool` em vez de `database/sql`
- `pgx.ErrNoRows` → 404
- `RETURNING *` no INSERT para popular id e timestamps

**JOINs e relacionamentos:**
- `LEFT JOIN` para incluir contatos sem telefones
- Resultado SQL "achatado" → agregação no Go
- Tipos `pgtype.Int4`/`pgtype.Text` para colunas que viram nulas pelo JOIN
- `phones: []` vs `phones: null` no JSON (slice nil vs slice vazia)
- N+1 problem: por que o JOIN é melhor que listar contatos + 1 query por contato

Para os conceitos de JOIN, veja:

- Slides da aula: JOINs e relacionamentos 1:N: LEFT JOIN, agregação no Go (no SIGAA)
- `ex04/README.md` da Lista 4

## Relação com a Lista 4

Esta pasta corresponde ao **Ex01 + Ex04 da Lista 4** resolvidos em aula. Os outros dois exercícios da Lista 4 são extensões individuais:

- **Ex01** ✅ resolvido aqui — sqlc básico
- **Ex02** — Repository pattern (interface desacoplando handler do sqlc)
- **Ex03** — Filtros, paginação e parâmetros nomeados sqlc
- **Ex04** ✅ resolvido aqui — JOIN 1:N e agregação no Go

Cada exercício tem README próprio com explicação detalhada. Os Ex02 e Ex03 são o trabalho **individual** da Sprint 2.

## Referências

- sqlc: https://docs.sqlc.dev
  - [Datatypes — PostgreSQL → Go](https://docs.sqlc.dev/en/stable/reference/datatypes.html)
- pgx: https://github.com/jackc/pgx
  - [`pgtype` package](https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype)
- PostgreSQL
  - [Foreign Keys (tutorial)](https://www.postgresql.org/docs/current/tutorial-fk.html)
  - [Table Expressions / JOINs](https://www.postgresql.org/docs/current/queries-table-expressions.html#QUERIES-FROM)
- Chi router: https://go-chi.io
- SvelteKit: https://kit.svelte.dev
- Problem Details (RFC 7807): https://datatracker.ietf.org/doc/html/rfc7807
