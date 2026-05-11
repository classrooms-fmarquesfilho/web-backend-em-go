# Frontend — Webapp de demonstração

Webapp SvelteKit que consome a API Go (na pasta `../backend/`). Cobre o material das aulas da Sprint 2: CRUD com sqlc e JOINs com relacionamentos 1:N.

## Telas

| Rota | O que faz |
|------|-----------|
| `/` | Lista de contatos + formulário de criação + remoção |
| `/contacts/{id}` | Detalhe do contato com gestão de telefones (criar/apagar) |
| `/contacts-with-phones` | Visão agregada de **todos** os contatos com seus telefones — consome o endpoint `/contacts-with-phones` do backend (LEFT JOIN + agregação no Go) |

A rota `/contacts-with-phones` é mostrado na aula de JOINs. Ela mostra:

- Estatísticas: N contatos, M telefones totais, **1 requisição HTTP, 1 query SQL**
- Tempo de resposta em ms
- Botão para inspecionar o JSON cru retornado pelo backend
- Cards visuais para cada contato com seus telefones aninhados

## Como rodar

Esta pasta é parte do projeto integrado `07-postgres-sqlc/`. A forma mais simples é subir tudo via `docker-compose` lá da raiz:

```bash
cd ..                      # vá para 07-postgres-sqlc/
docker compose up --build
```

E abra http://localhost:5173.

### Standalone (sem docker compose)

```bash
# Requer: API Go rodando em http://localhost:8080 (ver ../backend/)
npm install
npm run dev
```

## Estrutura

```
src/
  app.html
  lib/
    api.js                                          ← cliente HTTP da API
  routes/
    styles.css                                      ← tokens de cor e estilos base
    +layout.svelte                                  ← header + footer
    +page.svelte                                    ← lista + criação de contatos
    contacts/[id]/+page.svelte                      ← detalhe + telefones
    contacts-with-phones/+page.svelte               ← visão agregada (JOIN)
vite.config.js                                      ← proxy /api → backend
```

## Cliente da API

`src/lib/api.js` exporta funções tipadas (via JSDoc) para cada endpoint:

| Função | Endpoint |
|--------|----------|
| `listContacts()` | `GET /contacts` |
| `getContact(id)` | `GET /contacts/{id}` |
| `createContact({name, email})` | `POST /contacts` |
| `deleteContact(id)` | `DELETE /contacts/{id}` |
| `listPhones(contactId)` | `GET /contacts/{id}/phones` |
| `createPhone(contactId, {label, number})` | `POST /contacts/{id}/phones` |
| `deletePhone(contactId, phoneId)` | `DELETE /contacts/{contactId}/phones/{phoneId}` |
| `listContactsWithPhones()` | `GET /contacts-with-phones` |

Os erros do backend (Problem Details, RFC 7807) são transformados em `Error` com a mensagem do `detail` para serem propagados naturalmente.

## Proxy do Vite

`vite.config.js` redireciona `/api/*` (do front) para o backend. Em duas situações:

- **Dev local manual** (`API_TARGET` não definido): `http://localhost:8080`
- **Dentro do docker-compose** (`API_TARGET=http://backend:8080`): aponta para o serviço `backend` na rede interna

Isso elimina problemas de CORS em dev e mantém o código do frontend agnóstico de onde o backend está hospedado.

## Stack

- SvelteKit 2.15 + Svelte 5 + Vite 6
- `@sveltejs/adapter-node` (build para Node)
- Sem dependências de UI — CSS direto com tokens em `routes/styles.css`

A escolha foi por uma stack pequena e direta, com foco em demonstrar o consumo da API. Não é um app de produção.
