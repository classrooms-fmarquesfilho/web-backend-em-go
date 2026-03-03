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

<!-- Slide 1 — Título -->

# O pacote net/http — Parte 1
## Handlers, httptest e Query Params

### Curso Web Backend — Aula 02a

<br>

*Pré-requisitos do Curso Básico:* &nbsp;
<span class="tag">Hello World</span> &nbsp;
<span class="tag">Structs · Métodos · Interfaces</span> &nbsp;
<span class="tag">Injeção de Dependência</span>

---

<!-- Slide 2 — Arquitetura HTTP -->

# Arquitetura Web: Cliente-Servidor

```
   Cliente                         Servidor (Go)
     │── GET /celsius?valor=100 ──────▶│
     │                                 │  handler processa
     │◀── 200 OK  "212°F" ────────────│
```

<div class="columns">
<div class="col">

**Requisição**
- **Método**: `GET` `POST` `PUT` `DELETE`
- **Caminho**: `/ping`, `/celsius`
- **Query params**: `?valor=100`
- **Headers**: `Content-Type`, `Authorization`
- **Body**: dados JSON (POST / PUT)

</div>
<div class="col">

**Resposta**
- `200` — OK
- `400` — requisição inválida
- `404` — não encontrado
- `500` — erro interno

</div>
</div>

---

<!-- Slide 3 — Métodos HTTP -->

# Métodos HTTP: semântica e intenção

Cada método carrega uma **intenção** — não é convenção, é contrato.

| Método | Intenção | Seguro? | Idempotente? |
|--------|----------|:-------:|:------------:|
| `GET` | Buscar — sem efeito colateral | ✓ | ✓ |
| `POST` | Criar ou processar | ✗ | ✗ |
| `PUT` | Substituir completamente | ✗ | ✓ |
| `DELETE` | Remover | ✗ | ✓ |

<br>

> **Seguro**: não modifica estado no servidor — caches e navegadores podem repetir livremente.
> **Idempotente**: repetir N vezes tem o mesmo efeito que fazer uma vez.

---

<!-- Slide 4 — http.Handler e io.Writer -->

# `http.Handler`, `ResponseWriter` e `io.Writer`

<div class="columns">
<div class="col">

```go
// A interface central
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Na prática: HandleFunc
http.HandleFunc("/ping",
  func(w http.ResponseWriter,
       r *http.Request) {
    // w = onde escrever a resposta
    // r = a requisição recebida
  })
```

</div>
<div class="col">

```go
type ResponseWriter interface {
    Header() http.Header
    Write([]byte) (int, error) // ← io.Writer!
    WriteHeader(statusCode int)
}
```

> `ResponseWriter` satisfaz `io.Writer`.
> `fmt.Fprintf(w, ...)` funciona —
> mesma mecânica de
> **Injeção de Dependência**.

</div>
</div>

---

<!-- Slide 5 — httptest -->

# Testando handlers sem servidor: `httptest`

```go
func TestMeuHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/ping", nil)  // requisição simulada
    rec := httptest.NewRecorder()                     // captura a resposta

    meuHandler(rec, req)                              // chama direto — sem servidor!

    if rec.Body.String() != "pong" { ... }
    if rec.Code != 200 { ... }
}
```

| | O que é | O que inspecionar |
|---|---|---|
| `httptest.NewRequest()` | Requisição simulada | método, URL, body |
| `httptest.NewRecorder()` | Implementa `ResponseWriter` | `.Body` · `.Code` · `.Header()` |

---

<!-- Slide 6 — O que vamos construir -->

# O que vamos construir neste vídeo

Dois handlers com TDD — do mais simples ao que lê dados da URL:

| # | Endpoint | Conceito central |
|---|----------|-----------------|
| 1 | `GET /ping` → `pong` | handler mínimo · status code · `t.Helper()` |
| 2 | `GET /celsius?valor=100` → `212°F` | query params · `r.URL.Query()` · função pura |

<br>

<span class="pill-red">RED</span> &nbsp; escrever o teste primeiro, ver falhar &nbsp;&nbsp;
<span class="pill-green">GREEN</span> &nbsp; implementar o mínimo para passar &nbsp;&nbsp;
<span class="pill-blue">REFACTOR</span> &nbsp; melhorar sem quebrar

<br>

> No próximo vídeo: **maps em Go**, `ServeMux`, múltiplas rotas e JSON.