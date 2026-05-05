---
marp: true
theme: default
paginate: true
backgroundColor: #ffffff
color: #1a1a2e
style: |
  section {
    font-family: 'Calibri', sans-serif;
    padding: 36px 48px;
    font-size: 1.3em;
  }
  h1 {
    font-family: 'Consolas', monospace;
    color: #1a56db;
    font-size: 1.6em;
    margin-bottom: 0.3em;
    border-bottom: 2px solid #e5e7eb;
    padding-bottom: 0.2em;
  }
  h2 {
    font-family: 'Consolas', monospace;
    color: #374151;
    font-size: 1.2em;
    margin-bottom: 0.25em;
  }
  h3 { color: #6b7280; font-size: 0.9em; margin: 0.2em 0; }
  strong { color: #b45309; }
  em { color: #6b7280; }
  code {
    font-family: 'Consolas', monospace;
    background: #e5e7eb;
    color: #1e3a5f;
    padding: 0.08em 0.3em;
    border-radius: 3px;
    font-size: 1.00em;
  }
  pre {
    background: #f3f4f6 !important;
    border: 1px solid #d1d5db;
    border-left: 3px solid #1a56db;
    border-radius: 6px;
    padding: 0.7em 1em;
    margin: 0.4em 0;
  }
  pre code {
    background: transparent;
    color: #1e3a5f;
    font-size: 0.9em;
    padding: 0;
    line-height: 1.5;
  }
  table { font-size: 0.95em; width: 100%; border-collapse: collapse; }
  th {
    background: #e5e7eb;
    color: #1a56db;
    font-family: 'Consolas', monospace;
    padding: 0.35em 0.7em;
    border: 1px solid #d1d5db;
  }
  td { background: #ffffff; padding: 0.28em 0.7em; border: 1px solid #d1d5db; color: #1a1a2e; }
  tr:nth-child(even) td { background: #f9fafb; }
  ul { margin: 0.25em 0; padding-left: 1.3em; }
  li { margin: 0.18em 0; font-size: 0.88em; line-height: 1.4; }
  blockquote {
    border-left: 3px solid #1a56db;
    background: #eff6ff;
    padding: 0.4em 0.9em;
    margin: 0.5em 0;
    font-style: normal;
    color: #1e3a5f;
    border-radius: 0 5px 5px 0;
    font-size: 1.00em;
  }
  .columns { display: flex; gap: 1.8em; }
  .col { flex: 1; }
  .pill-red   { display:inline-block; background:#fee2e2; border:1.5px solid #dc2626; color:#dc2626; font-family:'Consolas',monospace; font-weight:bold; font-size:0.85em; padding:0.12em 0.6em; border-radius:20px; }
  .pill-green { display:inline-block; background:#dcfce7; border:1.5px solid #16a34a; color:#16a34a; font-family:'Consolas',monospace; font-weight:bold; font-size:0.85em; padding:0.12em 0.6em; border-radius:20px; }
  .pill-blue  { display:inline-block; background:#dbeafe; border:1.5px solid #1a56db; color:#1a56db; font-family:'Consolas',monospace; font-weight:bold; font-size:0.85em; padding:0.12em 0.6em; border-radius:20px; }
  .tag { display:inline-block; background:#f3f4f6; border:1px solid #d1d5db; color:#374151; font-size:0.85em; padding:0.1em 0.5em; border-radius:4px; font-family:'Consolas',monospace; }

---

# PostgreSQL + sqlc
## Persistência type-safe para APIs Go

**DIM0547 — Sprint 2 (28/04)**

**Prof. Fernando** · UFRN · 2026.1

---

# O problema

```go
var contacts = map[string]Contact{}  // Sprint 1
```

- Reinicia o servidor → **perde tudo**
- Não escala para múltiplas instâncias
- Sem transações, sem índices, sem integridade referencial

**Sprint 2**: trocar por **PostgreSQL** de verdade.

A pergunta é: **como acessar o banco em Go?**

---

# Três abordagens

| Abordagem | Exemplo | Prós | Contras |
|-----------|---------|------|---------|
| **ORM** | GORM | Produtividade rápida | Abstrações vazam, queries escondidas |
| **Raw SQL** | `database/sql` | Controle total | Repetitivo, sem type safety |
| **SQL-first codegen** | **sqlc** ✅ | SQL real + tipo seguro | Precisa saber SQL |

---

# Por que sqlc (e não GORM)

**GORM** (ORM):
```go
db.Where("email = ?", email).First(&contact)
// Que SQL isso gera? 🤷
// E se a query for lenta? Debug no escuro.
```

**sqlc** (SQL-first):
```sql
-- name: GetContactByEmail :one
SELECT * FROM contacts WHERE email = $1;
```
```go
contact, err := queries.GetContactByEmail(ctx, email)
// SQL visível, tipado, compile-time checked
```

> **ORM esconde SQL. Quando a query fica lenta, vocês vão debugar SQL de qualquer jeito.** Melhor começar sabendo SQL.

---

# Como sqlc funciona

```
1. Você escreve SQL  →  2. sqlc gera Go  →  3. Você usa Go tipado
```

```
db/schema/001.sql        sqlc generate        internal/db/
  CREATE TABLE ...    ──────────────────>     contacts.sql.go
                                              models.go
db/queries/contacts.sql                       db.go
  -- name: List :many
  SELECT * FROM ...
```

> Erros de SQL são pegos em **compile time**, não em runtime.

---

# Setup: sqlc.yaml

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries/"
    schema: "db/schema/"
    gen:
      go:
        package: "db"
        out: "internal/db"
```

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

---

# Passo 1: Definir o schema

```sql
-- db/schema/001_contacts.sql

CREATE TABLE contacts (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

> **SERIAL** = auto-incremento. **TIMESTAMPTZ** = timestamp com timezone. **UNIQUE** = banco garante unicidade, não o código Go.

---

# Passo 2: Escrever as queries

```sql
-- db/queries/contacts.sql

-- name: ListContacts :many
SELECT * FROM contacts ORDER BY created_at DESC;

-- name: GetContact :one
SELECT * FROM contacts WHERE id = $1;

-- name: CreateContact :one
INSERT INTO contacts (name, email)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1;
```

---

# Anotações sqlc

| Anotação | Retorno Go | Quando usar |
|----------|-----------|-------------|
| `:one` | `(Model, error)` | SELECT que retorna 1 registro |
| `:many` | `([]Model, error)` | SELECT que retorna N registros |
| `:exec` | `error` | DELETE, UPDATE sem retorno |
| `:execresult` | `(sql.Result, error)` | Quando precisa de RowsAffected |
| `:execrows` | `(int64, error)` | Atalho para RowsAffected |

**`RETURNING *`** → faz INSERT/UPDATE retornar o registro criado (use com `:one`)

---

# Passo 3: Gerar código

```bash
sqlc generate
```

Gera automaticamente em `internal/db/`:

```go
// models.go — structs
type Contact struct {
    ID        int32
    Name      string
    Email     string
    CreatedAt time.Time
}

// contacts.sql.go — funções
func (q *Queries) ListContacts(ctx context.Context) ([]Contact, error)
func (q *Queries) GetContact(ctx context.Context, id int32) (Contact, error)
func (q *Queries) CreateContact(ctx context.Context, arg CreateContactParams) (Contact, error)
func (q *Queries) DeleteContact(ctx context.Context, id int32) error
```

> **Struct gerada com tipos corretos**: `int32` para SERIAL, `time.Time` para TIMESTAMPTZ.

---

# Passo 4: Conectar ao banco

```go
import (
    "github.com/jackc/pgx/v5/pgxpool"
    "meuapp/internal/db"
)

func main() {
    ctx := context.Background()

    pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    queries := db.New(pool)

    router := handler.NewRouter(queries)
    http.ListenAndServe(":8080", router)
}
```

**DATABASE_URL**: `postgres://user:pass@localhost:5432/mydb?sslmode=disable`

---

# Passo 5: Usar nos handlers

```go
func (h *Handler) listContacts(w http.ResponseWriter, r *http.Request) {
    contacts, err := h.queries.ListContacts(r.Context())
    if err != nil {
        writeProblem(w, ProblemDetails{
            Status: 500, Title: "Internal Error",
            Detail: "failed to list contacts",
        })
        return
    }
    writeJSON(w, http.StatusOK, contacts)
}
```

**Diferença do map**:
- Antes: `list := make([]Contact, 0); for _, c := range a.contacts { ... }`
- Depois: `contacts, err := h.queries.ListContacts(r.Context())`

> **Mais simples, mais seguro, persistente.**

---

# Migrations: versionando o schema

**Problema**: como aplicar mudanças no banco de forma **reprodutível e versionada**?

**Migrations** = scripts SQL numerados que transformam o schema:

```
db/schema/
├── 001_contacts.sql        # CREATE TABLE contacts
├── 002_add_phone.sql       # ALTER TABLE contacts ADD COLUMN phone TEXT
└── 003_create_orders.sql   # CREATE TABLE orders
```

Cada migration roda **uma vez**, em ordem. O banco lembra quais já foram aplicadas.

---

# Aplicando migrations (simples)

**Opção mínima** (sem ferramenta):

```bash
psql -U user -d mydb -f db/schema/001_contacts.sql
```

**Com golang-migrate** (recomendado para CI):

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate -path db/schema -database "$DATABASE_URL" up
```

**Com Atlas** (mais poderoso):

```bash
atlas migrate apply --url "$DATABASE_URL"
```

> Para a Sprint 2, **qualquer uma das três opções** é aceitável. O que importa é que o schema esteja versionado em arquivo.

---

# Docker: PostgreSQL local

Para desenvolver, **não instale PostgreSQL na máquina**. Use Docker:

```bash
docker run -d \
  --name postgres-dev \
  -p 5432:5432 \
  -e POSTGRES_USER=dev \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=myapp \
  postgres:15
```

Ou com **docker-compose.yml** (Sprint 2 pede):

```yaml
services:
  db:
    image: postgres:15
    ports: ["5432:5432"]
    environment:
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: myapp
```

```bash
docker compose up -d
```

---

# Estrutura do projeto Sprint 2

```
meu-projeto/
├── cmd/api/main.go
├── internal/
│   ├── db/                        ← NOVO (gerado pelo sqlc)
│   │   ├── contacts.sql.go
│   │   ├── db.go
│   │   └── models.go
│   ├── handler/
│   │   └── contacts.go            ← usa db.Queries
│   └── middleware/
├── db/
│   ├── schema/001_contacts.sql    ← NOVO (migrations)
│   └── queries/contacts.sql       ← NOVO (queries SQL)
├── sqlc.yaml                      ← NOVO
├── docker-compose.yml             ← NOVO
├── go.mod
└── .github/workflows/ci.yml
```

---

# O que mudou do Sprint 1 para o Sprint 2

| Sprint 1 | Sprint 2 |
|----------|----------|
| `map[string]Contact` | PostgreSQL + sqlc |
| Dados em memória | Dados persistentes |
| Sem schema | Schema versionado |
| Sem Docker | docker-compose |
| Tudo em `handler/` | Separação handler ↔ db |

> Sprint 3 vai separar ainda mais: **Clean Architecture** com camadas service/repository.

---

# Resumo

1. **sqlc** gera Go tipado a partir de SQL puro
2. **Schema** versionado em arquivos SQL (`db/schema/`)
3. **Queries** anotadas com `:one`, `:many`, `:exec`
4. **`sqlc generate`** cria structs + funções automaticamente
5. **PostgreSQL via Docker** para desenvolvimento local
6. **Migrations** para aplicar schema de forma reprodutível

**Para fazer agora**:
- Sprint 1 (quem não entregou): **até quinta 30/04**
- Sprint 2: conectar API ao PostgreSQL — entrega **08/05 (sex)**
- Lista 4 (sqlc + Clean Architecture): publicada quinta