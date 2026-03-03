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

<!-- Slide 1 — GitHub Classroom: o fluxo -->

# GitHub Classroom: o fluxo da Lista 1

**O que acontece quando você aceita o assignment:**

1. Repositório privado criado automaticamente: `lista-01-SEU-USUARIO`
2. Código base copiado (ex01–ex04 com `main.go` e `main_test.go`)
3. Pipeline de autograding configurado

**O ciclo de entrega:**

```
Discord #exercicios         GitHub                       Nota
─────────────────           ──────────────               ────
Link do assignment   →   git clone   →   implementar
                                             ↓
                         aba Actions  ←   git push
                         (resultado)
```

> Podem fazer **quantos pushes quiserem** antes do prazo.
> Cada push roda o autograding do zero.
> **Prazo**: 20/03/2026 às 23:59 — vale o horário do último commit.

---

<!-- Slide 2 — A contradição do main.go -->

# A contradição do `main.go`

O `main.go` de cada exercício mostra exemplos com `curl`:

```
// - Porta: 8080
// Exemplo:
//   curl -X POST http://localhost:8080/echo -d '{"message":"hello"}'
```

Instinto natural → `http.ListenAndServe(":8080", nil)` no `main()`.

<div class="warn">
⚠️ <strong>Não façam isso de forma incondicional.</strong>
O autograding vai esperar 60 segundos e reprovar com zero.
</div>

**Por quê?** O arquivo de teste tem:

```go
func init() {
    main()  // executada antes de qualquer teste
}
```

Se `main()` chamar `ListenAndServe`, o processo **bloqueia** — os testes nunca rodam.

<div class="ok">
✅ Os exemplos de <code>curl</code> são para <strong>teste manual no terminal</strong>.
Para o autograding, o <code>main()</code> só precisa registrar os handlers.
</div>

---

<!-- Slide 3 — Como os testes funcionam -->

# Como os testes realmente funcionam

<div class="columns">
<div class="col">

**O que o teste faz:**

```go
// 1. init() chama main() → handlers registrados
func init() { main() }

// 2. teste cria requisição em memória
req := httptest.NewRequest(
    http.MethodPost, "/echo",
    bytes.NewBufferString(input))

// 3. chama o mux global diretamente
rec := httptest.NewRecorder()
http.DefaultServeMux.ServeHTTP(rec, req)

// 4. inspeciona a resposta
res := rec.Result()
defer res.Body.Close()
```

</div>
<div class="col">

**O que o `main()` deve fazer:**

```go
func main() {
    // registrar no DefaultServeMux
    http.HandleFunc("/echo", handlerEcho)

    // flag de guarda para teste manual
    if len(os.Args) > 1 &&
       os.Args[1] == "-serve" {
        log.Fatal(
            http.ListenAndServe(":8080", nil))
    }
}
```

- `go test ./...` → sem argumento → sem servidor
- `go run main.go -serve` → com servidor

</div>
</div>

> Use `http.HandleFunc` (não `http.NewServeMux()`) — os testes chamam `http.DefaultServeMux` diretamente.

---

<!-- Slide 4 — O pipeline: classroom.yml -->

# O pipeline: `.github/workflows/classroom.yml`

<div class="columns">
<div class="col">

**Triggers**
```yaml
on:
  - push            # automático em todo push
  - workflow_dispatch  # manual pela aba Actions
if: github.actor != 'github-classroom[bot]'
```

**Ambiente**
```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.22'  # ← CI usa 1.22, não 1.26
```

**Cada exercício é um step independente**
```yaml
- name: Exercício 4 - JSON Echo
  id: ex04
  with:
    command: cd ex04 && go test -v ./...
    timeout: 60     # segundos — ListenAndServe trava aqui
    max-score: 25
```

</div>
<div class="col">

**Reporter: calcula a nota final**
```yaml
- name: Autograding Reporter
  env:
    EX01_RESULTS: "${{steps.ex01.outputs.result}}"
    EX02_RESULTS: "${{steps.ex02.outputs.result}}"
    EX03_RESULTS: "${{steps.ex03.outputs.result}}"
    EX04_RESULTS: "${{steps.ex04.outputs.result}}"
  with:
    runners: ex01,ex02,ex03,ex04
```

**Pontuação total**: 0–100
(4 exercícios × 25 pontos)

> Steps são **independentes**: se ex01 travar,
> ex02, ex03 e ex04 ainda rodam.

</div>
</div>

---

<!-- Slide 5 — Lendo os resultados na aba Actions -->

# Lendo os resultados na aba Actions

**Onde encontrar**: repositório no GitHub → aba **Actions** → execução mais recente → job `run-autograding-tests`

| Step | Indica | Pontuação |
|------|--------|-----------|
| ✅ verde | todos os testes do exercício passaram | 25/25 |
| ❌ vermelho | pelo menos um teste falhou | 0/25 |

**Clique no step vermelho para ver o `go test -v` completo.**

**Erros mais comuns:**

<div class="warn">
<strong>Step vermelho sem output de testes — só "timeout"</strong><br>
Causa: <code>main()</code> tem <code>ListenAndServe</code> incondicional. O processo bloqueou por 60s.<br>
Correção: adicionar a flag <code>-serve</code> de guarda.
</div>

<div class="warn">
<strong><code>undefined: handlerEcho</code> (ou qualquer <code>undefined</code>)</strong><br>
Causa: função não implementada, ou nome diferente do que o teste usa.<br>
Correção: implementar a função com o nome correto.
</div>

---

<!-- Slide 6 — Checklist antes do push -->

# Checklist antes do push

<span class="num">1</span> `go test ./...` na raiz passa **sem nenhum vermelho**?
*Se falha local, falha no CI — não esperem o autograding pra descobrir.*

<span class="num">2</span> Nenhum `main()` tem `http.ListenAndServe` **sem a flag de guarda**?
*Revisar ex01, ex02, ex03 e ex04 — todos.*

<span class="num">3</span> O push foi para o branch **`main`**?
*O Classroom avalia o branch padrão.*

<span class="num">4</span> Depois do push, a aba **Actions** mostra ✅ no commit mais recente?
*Verde local não garante verde no CI — confirme sempre.*

<br>

> **Dúvidas?** Discord → `#duvidas-u1` — antes de 20/03 às 23:59.