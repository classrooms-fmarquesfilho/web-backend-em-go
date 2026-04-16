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

# Roteamento e Middleware
## net/http — Parte 2

**Prof. Fernando** · DIM0547 · 2026.1

*Go 1.22+: método no pattern, PathValue, padrão middleware*

---

# Roteamento no Go 1.22+

## Antes (Go < 1.22): tudo era genérico

```go
mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":    listarUsuarios(w, r)
    case "POST":   criarUsuario(w, r)
    default:       http.Error(w, "método não permitido", 405)
    }
})
```

## Agora (Go 1.22+): método no pattern

```go
mux.HandleFunc("GET /users", listarUsuarios)
mux.HandleFunc("POST /users", criarUsuario)
// 405 automático para outros métodos!
```

> Sem `switch`, sem `if`. O `ServeMux` resolve.

---

# Parâmetros de rota: `{nome}`

```go
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")   // ← extrai o parâmetro
    fmt.Fprintf(w, "Usuário: %s", id)
})
```

### Regras de precedência:

| Pattern | Match | Prioridade |
|---------|-------|-----------|
| `GET /users/{id}` | `/users/42` | Mais específico ganha |
| `GET /users/me` | `/users/me` | Literal > wildcard |
| `GET /` | Qualquer GET | Catch-all (menor prioridade) |

> **Dica**: `{nome...}` (com `...`) captura o **resto** do path: `/files/{path...}` casa com `/files/a/b/c`

---

# Exemplo completo: CRUD

```go
func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /contacts",       listar)
    mux.HandleFunc("POST /contacts",      criar)
    mux.HandleFunc("GET /contacts/{id}",  buscar)
    mux.HandleFunc("DELETE /contacts/{id}", remover)

    http.ListenAndServe(":8080", mux)
}

func buscar(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    contato, ok := store[id]
    if !ok {
        http.Error(w, "não encontrado", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(contato)
}
```

> Isso é o **Ex01 da Lista 2** — CRUD completo com ServeMux.

---

# O Padrão Middleware

## O que é?

Uma função que **envolve** um handler, executando código antes e/ou depois:

```
Request → [Logger] → [Auth] → [Handler] → Response
```

## A assinatura

```go
func MeuMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ANTES: log, auth, timing, etc.
        next.ServeHTTP(w, r)
        // DEPOIS: log de duração, cleanup, etc.
    })
}
```

> `http.HandlerFunc` é um **adapter** — converte uma função em `http.Handler`.

---

# Middleware de Log

```go
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        inicio := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s — %v", r.Method, r.URL.Path, time.Since(inicio))
    })
}
```

### Uso:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /ping", pingHandler)

// Envolve o mux inteiro:
http.ListenAndServe(":8080", Logger(mux))
```

> **Problema**: como capturar o **status code**? O `ResponseWriter` não expõe o status depois de `WriteHeader`. Solução: **status recorder** (Ex02 da Lista 2).

---

# Middleware de Auth (seletivo)

```go
func RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token != "Bearer secreto123" {
            http.Error(w, "não autorizado", http.StatusUnauthorized)
            return  // ← NÃO chama next — interrompe a cadeia
        }
        next.ServeHTTP(w, r)
    })
}
```

### Proteger apenas uma rota:

```go
mux.Handle("POST /admin/config", RequireAuth(
    http.HandlerFunc(adminConfig),
))
```

> Rotas GET continuam abertas. Só POST no admin precisa de token.

---

# Middleware de Recovery

```go
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("PANIC: %v", err)
                http.Error(w, `{"error":"internal server error"}`,
                    http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

**Sem recovery**: um `panic` derruba o servidor inteiro.
**Com recovery**: retorna 500, loga o erro, servidor continua rodando.

> Isso é o **Ex03 da Lista 2**.

---

# Compondo middlewares

Middlewares se compõem por aninhamento:

```go
handler := Logger(Recovery(Auth(mux)))
// Execução: Logger → Recovery → Auth → mux → Auth → Recovery → Logger
```

Função auxiliar para evitar aninhamento:

```go
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
    for i := len(mws) - 1; i >= 0; i-- {
        h = mws[i](h)
    }
    return h
}

// Uso:
handler := Chain(mux, Logger, Recovery, Auth)
```

---

# Context: passando dados entre middlewares

```go
type chaveRequestID struct{}

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := fmt.Sprintf("%d", time.Now().UnixNano())
        w.Header().Set("X-Request-Id", id)

        ctx := context.WithValue(r.Context(), chaveRequestID{}, id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// No handler downstream:
func meuHandler(w http.ResponseWriter, r *http.Request) {
    id := r.Context().Value(chaveRequestID{}).(string)
    // usa o id...
}
```

> Use **tipos privados** (`struct{}`) como chave — evita colisões entre pacotes.
> Isso é o **Ex04 da Lista 2**.

---

# 📝 Lista 2: Roteamento e Middleware

**Prazo**: 17/04 (sexta) — GitHub Classroom

| Ex | Tema | Conceito principal |
|----|------|--------------------|
| 01 | CRUD de contatos | `ServeMux` Go 1.22+, `PathValue`, JSON |
| 02 | Middleware de logging | `func(http.Handler) http.Handler`, status recorder |
| 03 | Middleware de recovery | `defer`, `recover()`, error handling graceful |
| 04 | Context e Request ID | `context.WithValue`, `r.WithContext`, headers |

### Dicas:
- Cada exercício exporta `NewRouter()` ou `NewApp()` — os testes chamam essas funções
- **Não** use `http.ListenAndServe` — testes usam `httptest` direto
- **Não** use `http.DefaultServeMux` — crie com `http.NewServeMux()`
- Leia os `*_test.go` — eles são a especificação

---

# Próximos passos

**Esta semana:**
- 📝 **Lista 2** (net/http + middleware) — publicada hoje, prazo **17/04 (sex)**
- Quinta 09/04: Aula presencial — **Introdução ao Chi** + JSON em Go

**Semana que vem (14-16/04):**
- Terça: Chi avançado + OpenAPI
- Quinta: 🔵 Acompanhamento online — **Sprint 1**
