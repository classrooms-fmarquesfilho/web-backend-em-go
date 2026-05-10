# Lista 04 — Webapp de demonstração

Webapp SvelteKit que consome a API Go do **ex01** da Lista 04.

CRUD de contatos: lista, criação, detalhe e remoção.

```
┌──────────────────┐         ┌──────────────────┐         ┌──────────────┐
│  SvelteKit (Vite)│         │ Vite Dev Proxy   │         │  API Go      │
│   localhost:5173 │  /api/* │   /api → :8080   │         │ localhost:   │
│                  │ ──────▶ │                  │ ──────▶ │   8080       │
└──────────────────┘         └──────────────────┘         └──────────────┘
                                                                  │
                                                                  ▼
                                                         ┌──────────────┐
                                                         │ PostgreSQL   │
                                                         └──────────────┘
```

## Pré-requisitos

- Node 20+
- A API do ex01 da Lista 04 rodando em `http://localhost:8080`
- PostgreSQL rodando (a API precisa)

## Como rodar

```bash
# 1. Suba o Postgres (se ainda não estiver rodando)
docker run -d --name lista04-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=lista04 \
  -p 5432:5432 \
  postgres:16-alpine

export DATABASE_URL="postgres://postgres:postgres@localhost:5432/lista04?sslmode=disable"

# 2. Aplique o schema e rode a API Go
cd ../lista-04-template/ex01
psql "$DATABASE_URL" -f db/schema/001_contacts.sql
go run ./cmd/api &

# 3. Em outro terminal, rode o webapp
cd lista-04-webapp
npm install
npm run dev
```

Abra `http://localhost:5173`.

## Estrutura

```
src/
  app.html
  lib/
    api.js                       ← cliente HTTP da API (fetch tipado via JSDoc)
  routes/
    +layout.svelte               ← header + footer
    +page.svelte                 ← lista + form de criação
    contacts/[id]/+page.svelte   ← detalhe + remoção
    styles.css                   ← variáveis e componentes base
vite.config.js                   ← proxy /api → http://localhost:8080
```

## Sobre o proxy

`vite.config.js` redireciona `/api/*` (do front) para `http://localhost:8080/*` (API Go), removendo o prefixo `/api`. Em desenvolvimento isso elimina problemas de CORS e permite chamar `fetch('/api/contacts')` direto.

Em produção, o webapp seria servido pelo mesmo gateway/reverse-proxy que serve a API, ou você ajustaria a constante `BASE` em `src/lib/api.js`.

## Por que SvelteKit?

Stack pequena, rápida de aprender, e suficiente para mostrar o ciclo completo de uma SPA consumindo uma API REST. Não há ambição de tornar isso um app de produção — é demonstração da Lista 04.
