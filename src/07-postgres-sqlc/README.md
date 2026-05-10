# 07 — PostgreSQL + sqlc

Material referente à aula de **05/05/2026** (Sprint 2) — primeira aula da unidade de persistência. Demonstra um CRUD de contatos com persistência real em PostgreSQL via [sqlc](https://sqlc.dev), e um frontend SvelteKit consumindo a API.

Esta pasta corresponde ao **Ex01 da Lista 4** resolvido em aula, mais um webapp pequeno de demonstração que consome a API.

## Arquitetura

```
┌──────────────┐  HTTP   ┌──────────────────┐ proxy /api ┌────────────┐  pgx  ┌────────────┐
│   Browser    │ ──────▶ │  SvelteKit (Vite)│ ─────────▶ │  Go + Chi  │ ────▶ │ PostgreSQL │
│  :5173       │         │  frontend        │            │  backend   │       │  db        │
└──────────────┘         │  :5173           │            │  :8080     │       │  :5432     │
                         └──────────────────┘            └────────────┘       └────────────┘
```

Cada caixa é um processo independente. Na aula os três são serviços do `docker-compose`, mas você pode rodar manualmente (ver opção B abaixo).

## Estrutura

```
07-postgres-sqlc/
├── README.md              ← este arquivo
├── docker-compose.yml     ← sobe os três serviços de uma vez
├── backend/               ← API Go (Ex01 da Lista 4 resolvido)
│   ├── SOLUCAO.md         ← raciocínio dos 4 TODOs
│   ├── handler/
│   ├── db/{schema,queries}/
│   ├── internal/db/       ← gerado pelo sqlc
│   └── ...
└── frontend/              ← Webapp SvelteKit
    ├── src/
    │   ├── routes/        ← páginas (+page.svelte, contacts/[id]/)
    │   └── lib/api.js     ← cliente HTTP da API
    └── ...
```

## Como rodar

### Opção A — Docker Compose

Subir os três serviços com um comando:

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
docker compose down -v    # apaga o banco também
```

### Opção B — Manual (3 terminais)

Útil para entender o que o compose faz por baixo, ou para mexer em um serviço só.

**Terminal 1 — Postgres:**

```bash
docker run -d --name lista04-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=lista04 \
  -p 5432:5432 postgres:16-alpine

# Aplicar o schema
psql "postgres://postgres:postgres@localhost:5432/lista04?sslmode=disable" \
  -f backend/db/schema/001_contacts.sql
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

Abre http://localhost:5173.

## Endpoints da API

```
GET    /contacts        → 200 + array JSON
POST   /contacts        → 201 + objeto criado (body: {"name":"...", "email":"..."})
GET    /contacts/{id}   → 200 ou 404
DELETE /contacts/{id}   → 204 ou 404
```

Erros usam **Problem Details** (RFC 7807): `Content-Type: application/problem+json`.

## Testar a API direto (sem o frontend)

```bash
# Criar
curl -X POST http://localhost:8080/contacts \
  -H 'Content-Type: application/json' \
  -d '{"name":"Maria","email":"maria@x.com"}'

# Listar
curl http://localhost:8080/contacts

# Detalhar
curl http://localhost:8080/contacts/1

# Apagar
curl -X DELETE http://localhost:8080/contacts/1
```

## Próximos passos (Lista 4)

Esta pasta cobre o **Ex01 da Lista 4** — a parte didática, resolvida em aula. Os outros três exercícios da Lista 4 estendem o mesmo CRUD:

- **Ex02** — Repository pattern: interface que desacopla handler do sqlc
- **Ex03** — Filtros, paginação, parâmetros nomeados sqlc
- **Ex04** — Tabela de telefones, JOIN 1:N, agregação no Go

Cada exercício tem README com explicação aprofundada. Comece pelo Ex02 só **depois** de fazer este Ex01 funcionar local.

## Referências

- sqlc: https://docs.sqlc.dev
- pgx: https://github.com/jackc/pgx
- Chi router: https://go-chi.io
- SvelteKit: https://kit.svelte.dev
- Problem Details (RFC 7807): https://datatracker.ietf.org/doc/html/rfc7807
