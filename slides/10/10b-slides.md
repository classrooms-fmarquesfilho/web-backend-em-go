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
  .ide { background:#fff7ed; border:2px solid #b45309; border-radius:8px; padding:0.55em 0.9em; margin:0.5em 0; font-family:'Consolas',monospace; color:#92400e; font-size:0.95em; }
  .ide strong { color:#b45309; }
  .ide-files { display:block; margin-top:0.3em; color:#1e3a5f; font-size:0.92em; }
  .center-title { text-align:center; }

---

# Segurança da API — resolvendo a Lista 5+6

**DIM0547 — Sprint 3 · Vídeo 8 (Parte 2 de Autenticação)**

**Prof. Fernando** · UFRN · 2026.1

> Segundo vídeo de autenticação. No anterior implementamos **login + JWT + middleware** (e o `ex01` com bcrypt). Aqui fechamos o que falta para a **Entrega Final**: **refresh tokens** e **correções de segurança (OWASP)** — resolvendo a **Lista 5+6** ao longo do caminho.

---

# De onde paramos — e como usar este vídeo

Aula 07 (Vídeo 7) entregou:

- `POST /login` confere a senha com **bcrypt** e devolve um **JWT** (`ex01` + início do `ex02`)
- `AuthMiddleware` valida assinatura, `exp` e `alg`, e injeta o `userID` no context

O que a **Entrega Final (23/06)** ainda exige hoje:

| Requisito | Onde resolvemos |
|---|---|
| Refresh token com rotação | `ex04` |
| **≥ 2 correções de segurança OWASP** | ownership (`ex03`) + rate limiting (`ex04`) |

> **Como o vídeo funciona**: a **teoria** fica nos slides; a **solução de cada exercício** eu mostro **no IDE**. Quando aparecer o quadro laranja <span class="tag">🖥️ SOLUÇÃO NO IDE</span>, estou alternando para o editor.

---

# Mapa: Lista 5+6 ↔ Entrega Final

A Lista 5+6 não será publicada como tarefa — resolvemos os exercícios **aqui** e o mesmo código sustenta a entrega do projeto.

| Exercício | Conceito | O que provê para a entrega |
|---|---|---|
| `ex01` | bcrypt + handler testável | base de senhas (Aula 07) |
| `ex02` | JWT (access token) | login + rotas protegidas |
| `ex03` | **ownership** (OWASP API3) | 1ª correção de segurança |
| `ex04` | refresh + **rate limiting** | rotação + 2ª correção (API4) |

> Conceitos do repositório: **1** = `ex01`; **2** = `ex02`; **3** = `+ ex03`; **4** = `+ ex04`.

---

# Parte 1 — Tokens de acesso e renovação

---

# Recap: o access token (JWT) da Aula 07

O login devolve um **access token** assinado (HS256), curto (15 min), que o cliente envia em cada request:

```
Authorization: Bearer eyJhbGc...
```

Três cuidados que já vimos:

- **`exp` curto** — um vazamento tem impacto limitado no tempo
- **`alg` validado no parse** — defesa contra o ataque `alg:none`
- **erro genérico** no login (`invalid credentials`) — não revela se o email existe

<div class="ide">🖥️ <strong>SOLUÇÃO NO IDE — ex02</strong> (fechamento)
<span class="ide-files">· <code>internal/auth/jwt.go</code> → <code>GenerateAccessToken</code> / <code>ValidateAccessToken</code>
· <code>handler/auth.go</code> → <code>Login</code> emitindo o token</span></div>

---

# O dilema do tempo de vida

```
exp curto (minutos)              exp longo (dias)
└─ vazamento dura pouco          └─ vazamento dura muito
└─ usuário relogga toda hora     └─ usuário fica logado
   (péssima UX)                     (boa UX, péssima segurança)
```

Segurança e boa experiência puxam para lados opostos. Um único token não resolve.

> A saída é **separar responsabilidades**: um token para *acessar* (curto, usado o tempo todo) e outro para *renovar* (longo, usado raramente). Cada um otimiza um objetivo.

---

# Dois tokens, dois papéis

| | Access token | Refresh token |
|---|---|---|
| **Função** | Provar identidade em cada request | Obter um novo access token |
| **Tempo de vida** | Curto — 15 min | Longo — 7 dias |
| **Usado** | Em toda requisição | Só quando o access expira |
| **É stateless?** | Sim — JWT autocontido | Não — guardado (em *hash*) no banco |

> O token que viaja muito (access) dura pouco. O que dura muito (refresh) quase não viaja — e, por estar no banco, **pode ser revogado**.

---

# Por que o refresh fica no banco

```sql
CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,           -- hash, nunca o token original
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ              -- NULL = ativo
);
```

> Modelagem 1:N da Sprint 2 — um usuário tem muitos refresh tokens (um por dispositivo). `revoked_at NULL` é o estado "ativo"; preencher a data é revogar. Guardamos o **hash** (SHA-256): vazou o banco, os tokens não funcionam.

---

# O fluxo: rotação e detecção de reuso

```
login            → access + refresh-1   (refresh-1 salvo como hash)
/refresh (rt-1)  → revoga rt-1, emite access + refresh-2   (ROTAÇÃO)
/refresh (rt-1)  → rt-1 já revogado → REUSO detectado → 401
                   + revoga TODOS os refresh do usuário
```

- **Rotação**: cada `/refresh` invalida o token usado e emite um novo. O cliente sempre tem o último.
- **Detecção de reuso**: se um refresh **já revogado** reaparece, é sinal de roubo → revoga toda a família e força novo login.

> A ordem importa: a checagem de **reuso vem antes da de expiração**. Um token revogado que volta não é "expirado", é suspeito.

---

# Logout: o que o refresh resolve

| Token | Após o logout |
|---|---|
| Access token | Ainda válido — mas expira em minutos |
| Refresh token | Revogado **imediatamente** — não renova mais |

O `/logout` revoga o refresh e devolve **204 sempre** — exista ou não o token (idempotente). Diferenciar vazaria se aquele token chegou a existir.

<div class="ide">🖥️ <strong>SOLUÇÃO NO IDE — ex04</strong> (refresh tokens)
<span class="ide-files">· <code>internal/auth/token.go</code> → <code>GenerateRefreshToken</code> (crypto/rand) / <code>HashRefreshToken</code> (SHA-256)
· <code>handler/auth.go</code> → <code>Login</code> (emite o par), <code>Refresh</code> (rotação + reuso), <code>Logout</code></span></div>

---

# Parte 2 — Autorização e Segurança (OWASP)

---

# Autenticação ≠ Autorização

| | Pergunta | Quem resolve |
|---|---|---|
| **Autenticação** | *Quem é você?* | login + JWT (Aula 07) |
| **Autorização** | *Você pode fazer isto?* | regra de negócio no handler |

O JWT prova **quem** é o usuário. Ele **não** decide o que esse usuário pode acessar — isso é com a gente, em cada rota.

```
401 Unauthorized → não sei quem você é      (falta/inválido o token)
403 Forbidden    → sei quem é, mas não pode  (sem permissão)
404 Not Found    → o recurso não existe… ou você não pode saber que existe
```

> A distinção 401/403/404 é o vocabulário de toda esta parte. Guarde o **404 com segundo sentido** — voltamos a ele no BOLA.

---

# OWASP API Security Top 10

A **OWASP** mantém uma lista dos riscos mais comuns. APIs têm a **sua própria** lista (separada da web tradicional), porque os ataques são diferentes: não há tela, e o cliente fala direto com os endpoints.

Hoje abordaremos as duas mais frequentes:

| Risco | Nome | Onde, no projeto |
|---|---|---|
| **API3:2023** | *Broken Object Level Authorization* (BOLA) | acesso às notas (`ex03`) |
| **API4:2023** | *Unrestricted Resource Consumption* | força bruta no login (`ex04`) |

> BOLA é, há anos, o **risco nº 1** em APIs. É também o mais simples de explorar — e, por isso, o mais cobrado em revisão.

---

# API3 — BOLA, explicado

**Object Level Authorization** = verificar, a cada acesso, se *aquele objeto específico* pertence a quem está pedindo.

O ataque (também chamado **IDOR**):

```
Ana está logada e acessa a própria nota:
    GET /notes/7   (Authorization: Bearer <token da Ana>)

Ana troca o id na URL:
    GET /notes/8   ← a nota 8 é de João
```

Se a API devolver a nota 8, **quebrou** o controle de autorização por objeto. O token da Ana é válido (autenticação OK), mas ela **não é dona** da nota 8 (autorização falhou).

> A falha não está no token — está em **esquecer de checar o dono**. Autenticar não é autorizar.

---

# API3 — a defesa: 404, não 403

A regra de ouro da autorização: recurso de outro usuário → **404**, nunca **403**.

```
403 → "existe, mas você não pode"   ← confirma que a nota 8 existe
404 → "não encontrado"               ← a existência fica indistinguível
```

<div class="ide">🖥️ <strong>SOLUÇÃO NO IDE — ex03</strong> (ownership)
<span class="ide-files">· <code>handler/notes.go</code> → <code>Get</code> / <code>Update</code> / <code>Delete</code> traduzindo <code>ErrNoteNotFound</code> em 404</span></div>

---

# API4 — limitação de taxa

APIs gastam recursos a cada request: CPU, banco, e-mail, SMS. Sem limite, um cliente abusivo consegue:

- **Força bruta** de senha no `/login` — milhares de tentativas por segundo
- **Negação de serviço** (DoS) — afogar o servidor de requisições
- Estouro de custo em serviços pagos (envio de SMS, etc.)

A defesa é **limitar a taxa** (rate limiting) por origem (IP):

```
até 5 req/s por IP, com folga (burst) de 10
acima disso  →  429 Too Many Requests  (+ header Retry-After)
```

> O `/login` é o alvo clássico de força bruta — por isso o rate limit entra exatamente nele. O algoritmo é o **token bucket**: cada IP tem um "balde" que reabastece a 5 fichas/s.

---

# API4 — rate limiting no login

O middleware fica **na frente** do handler de login: lê o IP de origem, consulta o limiter daquele IP e, se estourou, corta a requisição com **429** antes de tocar no banco.

- Um `rate.Limiter` por IP, guardados num `map` protegido por mutex
- `429` com `Retry-After: 1` e corpo JSON `{"error":"rate limit exceeded"}`

<div class="ide">🖥️ <strong>SOLUÇÃO NO IDE — ex04</strong> (rate limiting)
<span class="ide-files">· <code>middleware/rate_limit.go</code> → <code>Middleware</code> (SplitHostPort → getLimiter → Allow → 429)
· <code>cmd/api/main.go</code> → onde o limiter é plugado só no <code>/login</code></span></div>

---

# As correções de segurança, em um quadro

| Correção | OWASP / risco | Onde |
|---|---|---|
| Ownership por objeto (404, não 403) | **API3:2023** BOLA | `ex03` |
| Rate limiting no login | **API4:2023** | `ex04` |
| Mensagem genérica "invalid credentials" | enumeração de usuário | `ex02` |
| `alg` validado no parse | ataque `alg:none` | `ex02` |
| Senha em **bcrypt**, refresh em **SHA-256** | vazamento de banco | `ex01`/`ex04` |
| `JWT_SECRET` em *env var* | segredo no repositório | config |

> Para a entrega, **duas** correções já bastam.

---

# Referências

**Refresh tokens e sessão**

- [OWASP — Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [`golang-jwt/jwt`](https://github.com/golang-jwt/jwt) · [`golang.org/x/time/rate`](https://pkg.go.dev/golang.org/x/time/rate)

**OWASP API Security**

- [OWASP API Security Top 10 (2023)](https://owasp.org/API-Security/editions/2023/en/0x11-t10/) — visão geral
- [API3:2023 — Broken Object Level Authorization](https://owasp.org/API-Security/editions/2023/en/0xa1-broken-object-level-authorization/)
- [API4:2023 — Unrestricted Resource Consumption](https://owasp.org/API-Security/editions/2023/en/0xa4-unrestricted-resource-consumption/)