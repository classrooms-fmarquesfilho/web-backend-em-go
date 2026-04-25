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
# Chi Avançado, Validação e Middleware
## DIM0547 — Sprint 1

**Prof. Fernando** · UFRN · 2026.1


---

# Onde estamos

- **No vídeo anterior**: Chi básico — `chi.NewRouter()`, `r.Get/Post`, `chi.URLParam`, refatoração do ex01
- **Neste vídeo**: Chi avançado — grupos de rotas, middleware, validação, erros padronizados

Depois deste vídeo vocês têm tudo para:
- ✅ Fazer os 4 exercícios da **Lista 3**
- ✅ Entregar a **Sprint 1** do projeto


---

# PARTE 1 — Organizando rotas

---

# O problema: tudo em main.go

```go
func NewRouter() http.Handler {
    r := chi.NewRouter()
    r.Get("/api/v1/contacts", listContacts)
    r.Post("/api/v1/contacts", createContact)
    r.Get("/api/v1/contacts/{id}", getContact)
    r.Delete("/api/v1/contacts/{id}", deleteContact)
    r.Get("/api/v1/orders", listOrders)
    r.Post("/api/v1/orders", createOrder)
    r.Get("/api/v1/orders/{id}", getOrder)
    r.Get("/api/v1/products", listProducts)
    // ... mais 20 rotas ...
    return r
}
```

Com 30 endpoints, esse arquivo tem **500 linhas**. Impossível de manter.

---

# Grupos de rotas com r.Route

```go
r := chi.NewRouter()

r.Route("/api/v1", func(r chi.Router) {
    r.Get("/contacts", listContacts)
    r.Post("/contacts", createContact)
    r.Get("/contacts/{id}", getContact)
    r.Delete("/contacts/{id}", deleteContact)
})

r.Route("/api/v2", func(r chi.Router) {
    r.Get("/contacts", listContactsV2)  // formato diferente
})

r.Get("/health", healthCheck)  // fora dos grupos
```

**Prefixo declarado uma vez**, escopo claro. Versionamento de API com **zero repetição**.

> Isso é o **ex02** da Lista 3.

---

# Subrouters com r.Mount

Cada recurso em seu próprio pacote:

```go
// internal/handler/contacts.go
func ContactsRouter() chi.Router {
    r := chi.NewRouter()
    r.Get("/", listContacts)
    r.Post("/", createContact)
    r.Get("/{id}", getContact)
    r.Delete("/{id}", deleteContact)
    return r
}
```

No `main.go` — só composição:

```go
r := chi.NewRouter()
r.Mount("/api/v1/contacts", handler.ContactsRouter())
r.Mount("/api/v1/orders",   handler.OrdersRouter())
```

> Cada módulo cuida das suas rotas. O main só compõe.

---

# r.Route vs r.Mount: quando usar cada

<div class="two-col">

**`r.Route`** — rotas inline

```go
r.Route("/api/v1", func(r chi.Router) {
    r.Get("/contacts", list)
    r.Post("/contacts", create)
})
```

Bom para: **agrupar logicamente** dentro do mesmo arquivo

</div>

<div class="two-col">

**`r.Mount`** — router externo

```go
r.Mount("/api/v1/contacts",
    handler.ContactsRouter())
```

Bom para: **separar em pacotes** diferentes

</div>

Podem ser combinados:

```go
r.Route("/api/v1", func(r chi.Router) {
    r.Mount("/contacts", handler.ContactsRouter())
    r.Mount("/orders",   handler.OrdersRouter())
})
```

---

# PARTE 2 — Middleware

---

# Middleware: o padrão universal de Go

```go
func MeuMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ANTES do handler
        log.Println("request chegou:", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)  // chama o próximo

        // DEPOIS do handler
        log.Println("response enviada")
    })
}
```

**O padrão**: recebe `http.Handler`, retorna `http.Handler`.

> Esse padrão é da **stdlib de Go**, não do Chi. Funciona em qualquer lugar.
> É o **ex03** da Lista 3.

---

# Middleware global vs por grupo

```go
r := chi.NewRouter()

// ── Global: todas as rotas ──
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)

// ── Sem auth ──
r.Get("/health", healthHandler)
r.Get("/ping", pingHandler)

// ── Com auth (só /api) ──
r.Route("/api", func(r chi.Router) {
    r.Use(authMiddleware)       // ← só aqui dentro
    r.Get("/contacts", listContacts)
    r.Get("/orders", listOrders)
})
```

`/health` → Logger + Recoverer
`/api/contacts` → Logger + Recoverer **+ authMiddleware**

---

# Middleware stack: ordem importa

```go
r.Use(middleware.RequestID)   // 1º — gera ID
r.Use(middleware.Logger)      // 2º — loga (já com ID)
r.Use(middleware.Recoverer)   // 3º — captura panic
```

```
Request  →  RequestID  →  Logger  →  Recoverer  →  Handler
Response ←  RequestID  ←  Logger  ←  Recoverer  ←  Handler
```

> **RequestID antes de Logger** — senão o log não tem o ID.
> **Recoverer por último** — captura panics de tudo que vem antes.

---

# Middleware embutidos do Chi

`github.com/go-chi/chi/v5/middleware`

| Middleware | O que faz |
|------------|-----------|
| `Logger` | Log de método, path, status, duração |
| `Recoverer` | Captura panics → 500 em vez de crash |
| `RequestID` | Header `X-Request-Id` único por request |
| `RealIP` | IP real atrás de proxy/load balancer |
| `Throttle(n)` | Limita N requests simultâneos |
| `Timeout(d)` | Cancela request após duração |

```go
import chiware "github.com/go-chi/chi/v5/middleware"

r.Use(chiware.Logger)
r.Use(chiware.Recoverer)
```

> Na Sprint 1 vocês precisam de **pelo menos 1 middleware** (logging OU recovery).

---

# PARTE 3 — Validação

---

# O problema: validação manual

```go
func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)

    if req.Name == "" {
        http.Error(w, "name required", 400)
        return
    }
    if len(req.Name) < 2 {
        http.Error(w, "name too short", 400)
        return
    }
    if req.Email == "" {
        http.Error(w, "email required", 400)
        return
    }
    // ... 20 linhas de if/else ...
}
```

Repetitivo, propenso a erro, inconsistente. **Não escala.**

---

# Validação declarativa com go-playground/validator

```bash
go get github.com/go-playground/validator/v10
```

```go
type CreateUserRequest struct {
    Name     string `json:"name"     validate:"required,min=2,max=100"`
    Email    string `json:"email"    validate:"required,email"`
    Age      int    `json:"age"      validate:"required,min=18,max=120"`
    Password string `json:"password" validate:"required,min=8"`
}
```

**Tags `validate`** nos campos da struct. Sem `if/else`.

| Tag | Significado |
|-----|-------------|
| `required` | Campo obrigatório (não pode ser zero value) |
| `min=2` | Mínimo 2 (caracteres p/ string, valor p/ int) |
| `max=100` | Máximo 100 |
| `email` | Formato de email válido |

---

# Usando o validator no handler

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest

    // 1. Decodificar JSON
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // JSON inválido → 400
        writeProblem(w, ProblemDetails{
            Status: 400, Title: "Bad Request", Detail: "invalid json",
        })
        return
    }

    // 2. Validar
    if err := validate.Struct(&req); err != nil {
        // Validação falhou → 422 com lista de erros
        writeProblem(w, ProblemDetails{
            Status: 422, Title: "Validation Error",
            Detail: "Um ou mais campos são inválidos",
            Errors: extractFieldErrors(err),
        })
        return
    }

    // 3. Campos válidos — criar usuário
}
```

---

# Tags comuns do validator

| Tag | Exemplo | Valida |
|-----|---------|--------|
| `required` | `validate:"required"` | Campo não pode ser zero value |
| `email` | `validate:"email"` | Formato de email |
| `min` | `validate:"min=8"` | Mín. 8 chars (string) ou valor 8 (int) |
| `max` | `validate:"max=100"` | Máx. 100 |
| `len` | `validate:"len=11"` | Tamanho exato |
| `oneof` | `validate:"oneof=admin user"` | Um dos valores listados |
| `url` | `validate:"url"` | URL válida |
| `gte` | `validate:"gte=0"` | ≥ 0 (para números) |
| `e164` | `validate:"e164"` | Telefone internacional (+55...) |

Combinar com vírgula: `validate:"required,email"`

> Referência completa: [github.com/go-playground/validator](https://github.com/go-playground/validator)

---

# PARTE 4 — Erros padronizados

---

# O problema: cada API inventa seu formato de erro

```json
{"error": "not found"}
{"message": "invalid", "code": 422}
{"errors": ["name required", "email invalid"]}
{"success": false, "reason": "unauthorized"}
```

Cada API retorna erros de um jeito diferente. O cliente não sabe o que esperar.

**Solução**: **RFC 7807 — Problem Details for HTTP APIs**

Um formato padronizado que qualquer cliente sabe ler.

---

# Problem Details (RFC 7807 / 9457)

```json
{
  "type": "https://example.com/errors/validation",
  "title": "Validation Error",
  "status": 422,
  "detail": "Um ou mais campos são inválidos",
  "errors": [
    {"field": "email", "message": "must be a valid email"},
    {"field": "age", "message": "must be at least 18"}
  ]
}
```

**Content-Type**: `application/problem+json`

| Campo | Obrigatório | Significado |
|-------|-------------|-------------|
| `type` | sim | URI que identifica o tipo de erro |
| `title` | sim | Título legível do erro |
| `status` | sim | Código HTTP |
| `detail` | não | Descrição específica desta ocorrência |
| `errors` | não | Lista de erros por campo (extensão) |

---

# Implementando Problem Details em Go

```go
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type ProblemDetails struct {
    Type   string       `json:"type"`
    Title  string       `json:"title"`
    Status int          `json:"status"`
    Detail string       `json:"detail"`
    Errors []FieldError `json:"errors,omitempty"`
}

func writeProblem(w http.ResponseWriter, p ProblemDetails) {
    w.Header().Set("Content-Type", "application/problem+json")
    w.WriteHeader(p.Status)
    json.NewEncoder(w).Encode(p)
}
```

> `omitempty` em `Errors`: se não houver erros de campo, o array nem aparece no JSON.
> Isso é o **ex04** da Lista 3.

---

# Extraindo erros do validator para FieldError

```go
func extractFieldErrors(err error) []FieldError {
    var ve validator.ValidationErrors
    if !errors.As(err, &ve) {
        return nil
    }
    out := make([]FieldError, 0, len(ve))
    for _, fe := range ve {
        out = append(out, FieldError{
            Field:   strings.ToLower(fe.Field()),
            Message: messageFor(fe),
        })
    }
    return out
}

func messageFor(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required": return "is required"
    case "email":    return "must be a valid email"
    case "min":      return fmt.Sprintf("must be at least %s", fe.Param())
    case "max":      return fmt.Sprintf("must be at most %s", fe.Param())
    default:         return fmt.Sprintf("failed: %s", fe.Tag())
    }
}
```

---

# Padrão completo: decode → validate → respond

```go
func (a *app) createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeProblem(w, ProblemDetails{
            Type: "https://example.com/errors/bad-request",
            Title: "Bad Request", Status: 400,
            Detail: "invalid json body",
        })
        return
    }

    if err := a.validate.Struct(&req); err != nil {
        writeProblem(w, ProblemDetails{
            Type: "https://example.com/errors/validation",
            Title: "Validation Error", Status: 422,
            Detail: "Um ou mais campos são inválidos",
            Errors: extractFieldErrors(err),
        })
        return
    }

    user := User{ID: nextID(), Name: req.Name, Email: req.Email, Age: req.Age}
    a.users[user.ID] = user
    writeJSON(w, http.StatusCreated, user)
    // ⚠️ Senha NUNCA aparece na resposta (User não tem campo Password)
}
```

---

# PARTE 5 — Juntando tudo para a Sprint 1

---

# O que a Sprint 1 pede

| Requisito | Como resolver |
|-----------|---------------|
| API usando Chi como router | `chi.NewRouter()` + `r.Get/Post/Delete` |
| Pelo menos 4 endpoints | CRUD de 1 recurso já dá 4 |
| Pelo menos 1 middleware | `chiware.Logger` ou `chiware.Recoverer` |
| Testes automatizados verdes | `httptest` + `go test ./...` |
| GitHub Actions passando | CI em `.github/workflows/ci.yml` |
| Vídeo 8 min | Demonstrar API + tour pelo código |

**Avaliação**: funcionalidade (40%) + qualidade de código e testes (35%) + vídeo (25%)

---

# Estrutura de projeto recomendada

```
meu-projeto/
├── cmd/api/
│   └── main.go                    # Setup + ListenAndServe
├── internal/
│   ├── handler/
│   │   ├── contacts.go            # Router + handlers
│   │   └── contacts_test.go
│   ├── middleware/
│   │   ├── logging.go             # (ou use chiware.Logger)
│   │   └── recovery.go
│   └── model/
│       └── contact.go
├── go.mod
└── .github/workflows/ci.yml
```

> Não precisa ser exatamente assim. Mas **todo código em `main.go`** não é aceitável para a Sprint 1.

---

# Lista 3 → Sprint 1: conexão direta

| Exercício | O que pratica | Como aparece na Sprint 1 |
|-----------|---------------|--------------------------|
| **ex01** — Migrar para Chi | `chi.NewRouter`, `chi.URLParam` | Router do projeto |
| **ex02** — Grupos e subrouters | `r.Route`, `r.Mount` | Organização de rotas |
| **ex03** — Middleware | `r.Use`, escopo por grupo | Logging, recovery, auth |
| **ex04** — Validação + RFC 7807 | `validator`, ProblemDetails | Validação dos endpoints |

> Quem fizer a Lista 3 bem feita já tem **80% do código** da Sprint 1.

---

# Resumo do que vocês sabem agora

| Conceito | Slide | Lista 3 |
|----------|-------|---------|
| Grupos de rotas (`r.Route`) | ✅ | ex02 |
| Subrouters (`r.Mount`) | ✅ | ex02 |
| Middleware por grupo (`r.Use`) | ✅ | ex03 |
| Middleware customizado | ✅ | ex03 |
| `go-playground/validator` | ✅ | ex04 |
| Tags `validate:"required,email"` | ✅ | ex04 |
| Problem Details (RFC 7807) | ✅ | ex04 |
| Extrair erros do validator | ✅ | ex04 |
| CI com GitHub Actions | ✅ | Sprint 1 |
| Estrutura de projeto | ✅ | Sprint 1 |

**Próximo vídeo**: persistência com PostgreSQL e sqlc.