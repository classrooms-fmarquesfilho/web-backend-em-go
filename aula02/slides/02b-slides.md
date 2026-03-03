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

# O pacote net/http — Parte 2
## Maps, ServeMux e JSON

### Curso Web Backend — Aula 02b

<br>

*Continuação de:* &nbsp;
<span class="tag">net/http Parte 1</span> &nbsp;
<span class="tag">Handlers · httptest · Query Params</span>

---

<!-- Slide 2 — Maps em Go -->

# Maps em Go

Estrutura **chave → valor**. Tipo da chave e do valor declarados explicitamente.

<div class="columns">
<div class="col">

```go
// declaração com literal
capitais := map[string]string{
    "RN": "Natal",
    "SP": "São Paulo",
}

// acessar, adicionar, remover
fmt.Println(capitais["RN"]) // "Natal"
capitais["CE"] = "Fortaleza"
delete(capitais, "SP")

// chave inexistente → zero value (sem panic)
cidade := capitais["AM"]        // ""
cidade, existe := capitais["RN"] // verificar existência
```

</div>
<div class="col">

**Iterando com `for range`**

```go
for chave, valor := range capitais {
    fmt.Println(chave, "→", valor)
}

// só a chave
for chave := range capitais {
    fmt.Println(chave)
}

// só o valor
for _, valor := range capitais {
    fmt.Println(valor)
}
```

> ⚠️ **Ordem não garantida** — cada
> execução pode iterar em sequência
> diferente. Não assuma ordenação.

</div>
</div>

---

<!-- Slide 3 — interface{} -->

# `interface{}`: o tipo que aceita qualquer valor

Em Go, toda interface define métodos que um tipo deve implementar.
`interface{}` define **zero métodos** — qualquer tipo a satisfaz.

```go
var x interface{}
x = 42           // int
x = "natal"      // string
x = true         // bool
x = []int{1, 2}  // slice
```

> Use com moderação: você perde a verificação de tipos em tempo de compilação.

**Por que importa com JSON?**

JSON pode ter campos de tipos variados num mesmo objeto:

```json
{"online": true, "versao": "1.0", "porta": 8080}
```

`map[string]interface{}` representa isso naturalmente — chave sempre `string`, valor qualquer coisa.

---

<!-- Slide 4 — JSON e maps -->

# JSON e maps: serializar e deserializar

<div class="columns">
<div class="col">

**Go → JSON**
```go
dados := map[string]interface{}{
    "online": true,
    "versao": "1.0",
}
json.NewEncoder(w).Encode(dados)
// → {"online":true,"versao":"1.0"}
```

</div>
<div class="col">

**JSON → Go**
```go
var r map[string]interface{}
json.NewDecoder(req.Body).Decode(&r)

// ⚠️ números viram float64
r["porta"]  // 8080.0, não 8080
```

</div>
</div>

> Quando a estrutura é **conhecida**, prefira uma **struct tipada** — mais segura em tempo de compilação.

```go
type statusResposta struct {
    Online bool   `json:"online"`
    Versao string `json:"versao"`
}
```

---

<!-- Slide 5 — ServeMux -->

# `ServeMux`: roteando requisições

<div class="columns">
<div class="col">

```go
// crie seu próprio mux
mux := http.NewServeMux()

mux.HandleFunc("/ping",   handlerPing)
mux.HandleFunc("/status", handlerStatus)
mux.HandleFunc("/",       handlerRaiz) // catch-all

log.Fatal(http.ListenAndServe(":8080", mux))
```

</div>
<div class="col">

**Por que não o `DefaultServeMux`?**

O `DefaultServeMux` é global — não dá pra injetar no teste sem subir servidor.

Com mux próprio, no teste basta:

```go
mux := configurarRotas()
mux.ServeHTTP(rec, req) // direto!
```

> `/` é **catch-all**: rotas sem handler
> específico chegam aqui. Verifique
> `r.URL.Path` para retornar 404.

</div>
</div>

---

<!-- Slide 6 — struct vs map -->

# Struct tipada vs `map[string]interface{}`

| | `map[string]interface{}` | struct tipada |
|---|---|---|
| Estrutura conhecida? | Não precisa | Sim — campos declarados |
| Erros de campo | Tempo de execução | **Tempo de compilação** |
| Legibilidade | Chaves como strings | Campos com nome e tipo |
| Uso típico | JSON dinâmico / ad hoc | Entrada de dados, contratos de API |

---

<!-- Slide 6b — rec.Result e map mutation -->

# Dois padrões para os exercícios

**Lendo a resposta via `rec.Result()`**

```go
res := rec.Result()
defer res.Body.Close()          // Body é io.ReadCloser — feche sempre
body, _ := io.ReadAll(res.Body) // lê todos os bytes de uma vez
json.Unmarshal(body, &data)     // deserializa os bytes
```

> Diferente de `rec.Body.String()` — aqui você tem um `*http.Response` completo,
> com `.StatusCode`, `.Header` e `.Body`. Padrão comum em testes de lista.

**Receber JSON, adicionar campo, devolver**

```go
var dados map[string]interface{}
json.NewDecoder(r.Body).Decode(&dados)  // decodifica entrada

dados["received"] = true                // adiciona campo ao map

json.NewEncoder(w).Encode(dados)        // serializa de volta
```

---

<!-- Slide 7 — O que vamos construir -->

# O que vamos construir neste vídeo

Dois handlers com TDD — rotas múltiplas e processamento de JSON:

| # | Endpoint | Conceito central |
|---|----------|-----------------|
| 3 | `GET /status` → JSON · `GET /` · 404 | `ServeMux` · catch-all · `map[string]interface{}` |
| 4 | `POST /calcular` ← → JSON | decode · encode · struct tipada vs map |

<br>

<span class="pill-red">RED</span> &nbsp; escrever o teste primeiro, ver falhar &nbsp;&nbsp;
<span class="pill-green">GREEN</span> &nbsp; implementar o mínimo para passar &nbsp;&nbsp;
<span class="pill-blue">REFACTOR</span> &nbsp; extrair lógica pura do handler
