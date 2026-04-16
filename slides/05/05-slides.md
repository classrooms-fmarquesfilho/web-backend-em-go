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

# Chi — quando ServeMux não basta

## DIM0547 — Sprint 1 (16/04)

**Prof. Fernando** · UFRN · 2026.1

*Acabamos de resolver o ex01 com `net/http` puro. Agora vamos ver o que muda com Chi e por quê.*

---

# O que é Chi?

- **Lib pequena** (~3 mil linhas, 1 dependência)
- **Sobre `net/http`**: usa `http.Handler` da stdlib
- **Foco em uma coisa**: roteamento expressivo
- **Sem mágica**: o que você lê é o que executa

```bash
go get github.com/go-chi/chi/v5
```

> **Não é um framework.** Você não "entra no mundo do Chi" — você **adiciona** Chi a um projeto Go normal.

---

# Por que Chi (e não Gin, Echo, Fiber)?

| Aspecto | Chi | Gin/Echo/Fiber |
|---------|-----|----------------|
| Tipo do handler | `http.Handler` (padrão Go) | `gin.Context` (lock-in) |
| Middleware | `func(http.Handler) http.Handler` | API própria |
| Compatível com libs Go | ✅ Tudo funciona | ❌ Precisa wrapper |
| Curva de aprendizado | Baixa (é HTTP) | Média (é o framework) |
| Performance | Suficiente | Levemente mais rápido |

**Critério da escolha**: o que você aprende em Chi **transfere** para qualquer projeto Go. Em Gin, fica restrito ao Gin.

---

# Hello, Chi

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()

    r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Olá, mundo!"))
    })

    http.ListenAndServe(":8080", r)
}
```

> **Note**: o handler é `func(w http.ResponseWriter, r *http.Request)` — exatamente igual ao da stdlib. **Zero abstração nova.**

---

# Antes e depois: ServeMux vs Chi

<div class="two-col">

**ServeMux (Go 1.22+)**

```go
mux := http.NewServeMux()

mux.HandleFunc(
    "GET /contacts",
    listContacts,
)
mux.HandleFunc(
    "GET /contacts/{id}",
    getContact,
)

// id := r.PathValue("id")
```

</div>

<div class="two-col">

**Chi**

```go
r := chi.NewRouter()

r.Get("/contacts", listContacts)
r.Get("/contacts/{id}", getContact)

// id := chi.URLParam(r, "id")
```

</div>

**Diferenças mínimas até aqui**. Chi começa a ganhar quando o app cresce.

---

# Onde Chi ajuda: grupos de rotas

```go
r := chi.NewRouter()

r.Route("/api/v1", func(r chi.Router) {
    r.Get("/contacts", listContacts)
    r.Post("/contacts", createContact)
    r.Get("/contacts/{id}", getContact)
    r.Delete("/contacts/{id}", deleteContact)
})

r.Route("/api/v2", func(r chi.Router) {
    r.Get("/contacts", listContactsV2)
    // ...
})
```

**ServeMux**: você teria que repetir `/api/v1/...` em cada `HandleFunc`.
**Chi**: prefixo declarado uma vez, escopo claro.

---

# Subrouters e Mount

Quando o app cresce, você quer **separar arquivos**:

```go
// internal/api/contacts.go
func ContactsRouter() chi.Router {
    r := chi.NewRouter()
    r.Get("/", listContacts)
    r.Post("/", createContact)
    r.Get("/{id}", getContact)
    return r
}

// main.go
r := chi.NewRouter()
r.Mount("/api/v1/contacts", contacts.ContactsRouter())
r.Mount("/api/v1/orders",   orders.OrdersRouter())
```

> Cada módulo expõe seu próprio router. O `main` só compõe.

---

# Middleware: o padrão universal de Go

```go
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
```

**O padrão**: uma função que **recebe** `http.Handler` e **retorna** `http.Handler`.

> Esse padrão **não é do Chi**. É da stdlib. Funciona em ServeMux puro também.
> O que Chi adiciona é **conveniência** para aplicar middleware a grupos.

---

# Middleware por grupo

```go
r := chi.NewRouter()
r.Use(middleware.Logger)         // todos veem

r.Get("/health", healthHandler)  // só Logger

r.Route("/api", func(r chi.Router) {
    r.Use(middleware.RequestID)  // só /api tem
    r.Get("/info", infoHandler)
})

r.Route("/admin", func(r chi.Router) {
    r.Use(authMiddleware)        // só /admin exige auth
    r.Get("/users", listUsers)
})
```

**Escopo claro**. Sem confundir qual middleware se aplica onde.

---

# Middleware que vem com Chi

`github.com/go-chi/chi/v5/middleware`

| Middleware | O que faz |
|------------|-----------|
| `Logger` | Log de cada request |
| `Recoverer` | Captura panics (retorna 500 em vez de derrubar) |
| `RequestID` | Adiciona `X-Request-Id` único por request |
| `RealIP` | Detecta IP real atrás de proxies |
| `Throttle(n)` | Limita N requests simultâneos |
| `Timeout(d)` | Cancela request após duração |
| `Compress(n)` | Compressão gzip/deflate |

> Use `r.Use(middleware.Logger)` e ganhe logs estruturados em uma linha.

---

# 404 e 405 customizados

```go
r := chi.NewRouter()

r.NotFound(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "rota não encontrada",
    })
})

r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusMethodNotAllowed)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "método não permitido para esta rota",
    })
})
```

Em ServeMux puro: você teria que escrever um middleware que intercepta isso.

---

# Resumo: ServeMux vs Chi

| Caso de uso | ServeMux | Chi |
|-------------|----------|-----|
| API pequena (<10 endpoints) | ✅ Suficiente | ✅ Confortável |
| Versionamento `/api/v1` | 😐 Repetitivo | ✅ Natural |
| Middleware por grupo | ❌ Manual | ✅ Built-in |
| 404/405 customizado | 😐 Wrapper | ✅ Built-in |
| Subrouters em arquivos separados | 😐 Manual | ✅ `r.Mount` |
| Compatível com middleware de terceiros | ✅ | ✅ |

**Regra prática**: se o seu projeto vai crescer (e o seu **vai** crescer), comece com Chi.

---

# Para quem vem de Spring Boot

| Spring Boot | Chi |
|-------------|-----|
| `@RestController` | função que retorna `chi.Router` |
| `@RequestMapping("/api")` | `r.Route("/api", ...)` |
| `@GetMapping("/{id}")` | `r.Get("/{id}", handler)` |
| `@PathVariable` | `chi.URLParam(r, "id")` |
| `@RequestBody` | `json.NewDecoder(r.Body).Decode(&req)` |
| `@Valid` | `validator.Struct(&req)` (lib externa) |
| `Filter` / `Interceptor` | `func(http.Handler) http.Handler` |
| `@Autowired` | injeção manual no construtor |

**Os conceitos são os mesmos.** A sintaxe é mais explícita.

---

# Próximos passos

Agora vou refatorar **ao vivo** o ex01 que acabamos de fazer, migrando para Chi. Vocês vão ver:

1. ✅ Quanto código fica idêntico (handlers, storage, JSON)
2. ✅ Quanto fica mais simples (router setup)
3. ✅ Como adicionar middleware embutido em 1 linha
4. ✅ Como agrupar tudo em `/api/v1`

Depois apresento a **Lista 3** que cobre exatamente esses tópicos + validação.

---

# Referências

**Chi**
- [go-chi/chi no GitHub](https://github.com/go-chi/chi)
- [Documentação oficial](https://go-chi.io/)
- [Middleware embutido](https://github.com/go-chi/chi/tree/master/middleware)

**Padrões de middleware em Go**
- [Alex Edwards — Making and Using HTTP Middleware](https://www.alexedwards.net/blog/making-and-using-middleware)
- [Three Dots Labs — Common anti-patterns in Go web applications](https://threedots.tech/post/common-anti-patterns-in-go-web-applications/)

**Comparações**
- [Awesome Go HTTP routers](https://github.com/avelino/awesome-go#routers) — listagem completa
- [Why Go's net/http is enough for most apps](https://www.alexedwards.net/blog/which-go-router-should-i-use)