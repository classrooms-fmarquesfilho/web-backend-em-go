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

# Segurança da API

**DIM0547 — Sprint 3 · Vídeo 8 (Parte 2 de Autenticação)**

**Prof. Fernando** · UFRN · 2026.1

> Segundo vídeo de autenticação. No anterior implementamos **senhas com bcrypt** e o **login com JWT**. Aqui fechamos o que falta para a **Entrega Final**: **autorização**, **refresh tokens** e **proteção contra abuso** (através de limitação de taxa, ou _rate limiting_).

---

# De onde paramos — e como usar este vídeo

Aula 07 (Vídeo 7) entregou a **identidade**:

- `POST /login` confere a senha com **bcrypt** e devolve um **JWT**
- `AuthMiddleware` valida o token e injeta o `userID` no contexto

O que falta para a **Entrega Final (23/06)**:

| Tema | Exercício | Risco que resolve |
|---|---|---|
| Autorização por dono | `ex03` | OWASP **API1** (BOLA) |
| Refresh tokens | `ex04` | sessão sem reduzir segurança |
| Limitação de taxa | `ex05` | OWASP **API4** (abuso) |

> **Está tudo em memória.** Esta lista foca em **segurança**, não em persistência.

---

# Mapa: a Lista 5+6 em cinco passos

Cada exercício acrescenta **uma** ideia de segurança. A ordem importa: cada um se apoia no anterior.

| # | Exercício | Pergunta que responde |
|---|---|---|
| `ex01` | Senhas (bcrypt) | Como guardar senha sem guardar a senha? ✓ (Aula 07) |
| `ex02` | Access token (JWT) | Quem é o usuário em cada request? ✓ (Aula 07) |
| `ex03` | Autorização (ownership) | Este usuário **pode** acessar **este** objeto? |
| `ex04` | Refresh tokens | Como manter a sessão viva com segurança? |
| `ex05` | Limitação de taxa | Como impedir abuso de um endpoint? |

> Hoje resolvemos `ex03` → `ex04` → `ex05`, nessa ordem. `ex01` e `ex02` já estão prontos da aula passada — começamos com um recap rápido.

---

# Parte 1 — Identidade: quem é o usuário

---

# Recap: identidade já está resolvida

A Aula 07 fechou os dois primeiros exercícios:

- **`ex01` — bcrypt**: a senha nunca é guardada em texto puro; guardamos um *hash*. No login, `bcrypt.CompareHashAndPassword` confere se a senha bate com o *hash*.
- **`ex02` — JWT**: o login devolve um **access token** assinado (HS256), curto (15 min), enviado em cada request:

```
Authorization: Bearer eyJhbGc...
```

Três cuidados que já adotamos no `ex02`:

- **`exp` curto** — um vazamento tem impacto limitado no tempo
- **`alg` validado no parse** — defesa contra o ataque `alg:none`
- **erro genérico** no login (`invalid credentials`) — não revela se o e-mail existe

> Esta é a base. O access token prova **quem** é o usuário. Mas provar quem é **não** é o mesmo que decidir o que ele pode fazer — é aí que entra a Parte 2.

---

# Parte 2 — Autorização: o que o usuário pode

---

# Autenticação ≠ Autorização

| | Pergunta | Quem resolve |
|---|---|---|
| **Autenticação** | *Quem é você?* | login + JWT (Aula 07) |
| **Autorização** | *Você pode fazer **isto**?* | regra no handler |

O JWT prova a identidade. Ele **não** decide o que essa identidade pode acessar — isso é responsabilidade de cada rota.

```
401 Unauthorized → não sei quem você é      (falta/inválido o token)
403 Forbidden    → sei quem é, mas não pode  (sem permissão)
404 Not Found    → o recurso não existe… ou você não pode saber que existe
```

> Guarde o **404 com segundo sentido** — ele é a chave da defesa contra o BOLA, logo adiante.

---

# OWASP API Security Top 10

A **OWASP** mantém listas dos riscos de segurança mais comuns. APIs têm a **sua própria** lista, separada da web tradicional: não há tela, o cliente fala direto com os endpoints, e os ataques são outros.

Hoje abordaremos as duas mais frequentes:

| Risco | Nome | Onde, no projeto |
|---|---|---|
| **API1:2023** | *Broken Object Level Authorization* (BOLA) | acesso às notas (`ex03`) |
| **API4:2023** | *Unrestricted Resource Consumption* | abuso do login (`ex05`) |

> BOLA é, há anos, o **risco nº 1** em APIs — e o mais simples de explorar. É também o mais cobrado em revisão de código.

---

# API1 — BOLA, explicado

**Object Level Authorization** = verificar, a cada acesso, se *aquele objeto específico* pertence a quem está pedindo.

O ataque (também chamado **IDOR**):

```
Ana está logada e acessa a própria nota:
    GET /notes/7      (Bearer <token da Ana>)   → OK, é dela

Ana troca o id na URL:
    GET /notes/8      (Bearer <token da Ana>)   → a nota 8 é de João
```

Se a API devolver a nota 8, ela **quebrou** a autorização por objeto. O token da Ana é válido (**autenticação** OK), mas ela **não é dona** da nota 8 (**autorização** falhou).

> A falha não está no token — está em **esquecer de checar o dono**. Autenticar não é autorizar.

---

# API1 — a defesa: 404, não 403

A regra de ouro: recurso de outro usuário → **404**, nunca **403**.

```
403 → "existe, mas você não pode"   ← confirma que a nota 8 existe
404 → "não encontrado"               ← a existência fica indistinguível
```

Como garantir na prática:

- O *store* busca pela chave **e** confere o dono na mesma operação
- "não existe" e "existe, mas é de outro" devolvem o **mesmo** erro
- O handler traduz esse erro único para **404**

<div class="ide">🖥️ <strong>SOLUÇÃO NO IDE — ex03</strong> (ownership)
<span class="ide-files">· <code>internal/store/notes.go</code> → <code>GetForUser</code> / <code>UpdateForUser</code> / <code>DeleteForUser</code> (checagem <code>ok &amp;&amp; ownerID == userID</code>)
· <code>handler/notes.go</code> → traduz <code>ErrNoteNotFound</code> em 404</span></div>

---

# Parte 3 — Sessão: como manter o login ativo

---

# O dilema do tempo de vida

```
exp curto (minutos)                       exp longo (dias)
└─ vazamento dura pouco                   └─ vazamento dura muito
└─ usuário precisa autenticar toda hora   └─ usuário fica logado
   (péssima UX)                              (boa UX, péssima segurança)
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
| **É stateless?** | Sim — JWT autocontido | Não — guardado (em *hash*) no servidor |

> O token que viaja muito (access) dura pouco. O que dura muito (refresh) quase não viaja — e, por ser guardado no servidor, **pode ser revogado**.

---

# Por que o refresh é stateful

O access token é stateless: a assinatura basta para validá-lo. O refresh é o oposto — precisa de **estado** para podermos revogá-lo:

```go
type RefreshToken struct {
    Hash      string       // guardamos o HASH, nunca o token original
    UserID    uuid.UUID
    ExpiresAt time.Time
    RevokedAt *time.Time   // nil = ativo
}
```

> Um usuário tem **muitos** refresh tokens (um por dispositivo) — a relação 1:N que você já viu na Sprint 2, agora num *map* em memória. `RevokedAt == nil` é o estado "ativo".

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
· <code>internal/store/refresh.go</code> → <code>Insert</code> / <code>FindByHash</code> / <code>Revoke</code> / <code>RevokeAll</code>
· <code>handler/auth.go</code> → <code>Login</code> (emite o par), <code>Refresh</code> (rotação + reuso), <code>Logout</code></span></div>

---

# Parte 4 — Proteção contra abuso

---

# API4 — limitação de taxa

APIs gastam recursos a cada request: CPU, banco, e-mail, SMS. Sem limite, um cliente abusivo consegue:

- **Força bruta** de senha no `/login` — milhares de tentativas por segundo
- **Negação de serviço** (DoS) — afogar o servidor de requisições
- Estouro de custo em serviços pagos (envio de SMS, etc.)

A defesa é **limitar a taxa** por origem (IP):

```
até 5 req/s por IP, com folga (burst) de 10
acima disso  →  429 Too Many Requests  (+ header Retry-After)
```

> O `/login` é o alvo clássico de força bruta — por isso o rate limit entra exatamente nele. O algoritmo é o **token bucket**: cada IP tem um "balde" que reabastece a 5 fichas/s.

---

# API4 — rate limiting no login

O middleware fica **na frente** do handler de login: lê o IP de origem, consulta o *limiter* daquele IP e, se estourou, corta a requisição com **429** antes de passar a chamada adiante.

- Um `rate.Limiter` por IP, guardados num `map` 
- `429` com `Retry-After: 1` e corpo JSON `{"error":"rate limit exceeded"}`

<div class="ide">🖥️ <strong>SOLUÇÃO NO IDE — ex05</strong> (limitação de taxa)
<span class="ide-files">· <code>middleware/rate_limit.go</code> → <code>Middleware</code> (SplitHostPort → getLimiter → Allow → 429)
· <code>cmd/api/main.go</code> → onde o limiter é plugado só no <code>/login</code></span></div>

---

# As correções de segurança, em um quadro

| Correção | OWASP / risco | Onde |
|---|---|---|
| Autorização por objeto (404, não 403) | **API1:2023** BOLA | `ex03` |
| Limitação de taxa no login | **API4:2023** URC | `ex05` |
| Mensagem genérica "invalid credentials" | enumeração de usuário | `ex02` |
| `alg` validado no parse | ataque `alg:none` | `ex02` |
| Senha em **bcrypt**, refresh em **SHA-256** | vazamento de dados | `ex01`/`ex04` |
| `JWT_SECRET` em variável de ambiente | segredo no repositório | config |

> Para a Entrega Final, **duas** correções já bastam.

---

# Referências

**Identidade e sessão**

- [`golang-jwt/jwt`](https://github.com/golang-jwt/jwt) — access tokens
- [OWASP — Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [`golang.org/x/time/rate`](https://pkg.go.dev/golang.org/x/time/rate) — token bucket

**OWASP API Security**

- [OWASP API Security Top 10 (2023)](https://owasp.org/API-Security/editions/2023/en/0x11-t10/) — visão geral
- [API1:2023 — Broken Object Level Authorization](https://owasp.org/API-Security/editions/2023/en/0xa1-broken-object-level-authorization/)
- [API4:2023 — Unrestricted Resource Consumption](https://owasp.org/API-Security/editions/2023/en/0xa4-unrestricted-resource-consumption/)
