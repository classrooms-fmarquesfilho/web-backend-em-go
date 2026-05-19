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

# Autenticação com JWT

**DIM0547 — Sprint 3 · Aula 09**

**Prof. Fernando** · UFRN · 2026.1

---

# De onde paramos

Sprint 2: a API **persiste** dados de verdade — PostgreSQL, sqlc, JOINs 1:N.

| Já temos | Falta |
|---|---|
| CRUD persistente | Saber **quem** está chamando |
| Relacionamentos 1:N | Proteger rotas que não são públicas |
| Endpoints REST idiomáticos | Distinguir "qualquer um" de "o dono do recurso" |

Hoje a API ganha **identidade**: login que devolve um token, e rotas que só respondem para quem apresenta um token válido.

---

# Plano da aula

1. Autenticação vs autorização — não confundir <span class="tag">~5 min</span>
2. O problema: HTTP é sem estado <span class="tag">~10 min</span>
3. Anatomia de um JWT: três partes <span class="tag">~15 min</span>
4. Assinatura: por que ninguém forja o token <span class="tag">~10 min</span>
5. Fluxo de login + middleware de auth em Go <span class="tag">~20 min</span>
6. Erros comuns e o que cai na Sprint 3 <span class="tag">~10 min</span>

---

# Autenticação ≠ autorização

| | Pergunta | Exemplo |
|---|---|---|
| **Autenticação** | Quem é você? | "Sou a Maria, aqui está minha senha" |
| **Autorização** | Você pode fazer isso? | "Maria pode editar o contato 7?" |

A ordem importa: primeiro **autentica** (estabelece identidade), depois **autoriza** (verifica permissão).

> JWT resolve principalmente a **autenticação** — carrega a identidade de forma confiável. A autorização é o que **você** decide fazer com a identidade que o token carrega (papel, dono do recurso, escopo).

---

# O problema: HTTP não tem memória

Cada requisição HTTP é **independente**. O servidor não "lembra" que você fez login na requisição anterior.

```
POST /login      → "ok, você é a Maria"   (servidor esquece)
GET  /contacts   → "...e você é quem?"     (de novo do zero)
```

Duas famílias de solução:

| Abordagem | Onde fica o estado |
|---|---|
| **Sessão no servidor** | Servidor guarda; cliente só tem um `session_id` |
| **Token autocontido (JWT)** | O **token** carrega a identidade; servidor não guarda nada |

> JWT é **stateless**: o servidor não precisa de uma tabela de sessões. Isso escala horizontalmente sem esforço — qualquer instância valida o token sozinha.

---

# Sessão vs JWT: o trade-off

<div class="columns">
<div class="col">

**Sessão no servidor**

- Estado no servidor (memória ou Redis)
- Revogar é trivial: apaga a sessão
- Precisa de armazenamento compartilhado entre instâncias
- Cliente carrega só um ID opaco

</div>
<div class="col">

**JWT (stateless)**

- Estado no **token**, no cliente
- Escala sem armazenamento compartilhado
- Revogar é **difícil** — o token é válido até expirar
- Cliente carrega a identidade inteira

</div>
</div>

> Não existe escolha "certa" universal. JWT é bom em APIs distribuídas; o controle de sessão é uma boa quando revogação imediata é requisito. Na Sprint 3 usamos JWT — e na próxima aula vemos como **refresh tokens** atenuam o problema da revogação.

---

# JWT: o que é, concretamente

Um **JSON Web Token** é uma string de três partes separadas por ponto:

```
eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIiwibmFtZSI6Ik1hcmlhIn0.dQ7...
└──── header ────┘ └──────── payload ────────┘ └─ signature ─┘
```

```
HEADER.PAYLOAD.SIGNATURE
```

| Parte | Conteúdo | Codificação |
|---|---|---|
| **Header** | algoritmo e tipo | Base64URL de um JSON |
| **Payload** | as *claims* (a identidade) | Base64URL de um JSON |
| **Signature** | a prova de autenticidade | bytes assinados, Base64URL |

> Base64URL **não é criptografia**. Qualquer pessoa decodifica o payload e lê o conteúdo. Quem garante a integridade é a **assinatura** — não o sigilo.

---

# Header e Payload

```json
// HEADER
{ "alg": "HS256", "typ": "JWT" }
```

```json
// PAYLOAD — as claims
{
  "sub": "1",                  // subject: ID do usuário
  "name": "Maria",
  "role": "user",
  "iat": 1716120000,           // issued at
  "exp": 1716123600            // expiration
}
```

> **Nunca** coloque dado sensível no payload (senha, número de cartão). O payload é **legível por qualquer um**. Coloque apenas o necessário para identificar e autorizar: ID, papel, talvez escopo.

---

# Claims registradas (RFC 7519)

Algumas claims têm significado padronizado — bibliotecas as validam automaticamente:

| Claim | Significado | Por que importa |
|---|---|---|
| `sub` | Subject — quem o token representa | É o "quem é você" |
| `iss` | Issuer — quem emitiu | Validar a origem do token |
| `exp` | Expiration — quando expira | Token velho é rejeitado |
| `iat` | Issued at — quando foi emitido | Idade do token |
| `nbf` | Not before — válido a partir de | Token "agendado" |

> `exp` é a claim mais importante para a segurança prática. Sem ela, o token vale **para sempre** — um vazamento vira um problema permanente. Tokens de acesso devem ser **curtos** (minutos).

---

# A assinatura: o coração do JWT

```
signature = HMAC-SHA256(
    base64url(header) + "." + base64url(payload),
    SECRET
)
```

A assinatura é função de **três coisas**: o header, o payload e um **segredo** que só o servidor conhece.

| Alguém tenta... | Resultado |
|---|---|
| Mudar `"role": "user"` para `"admin"` | A assinatura não confere mais |
| Recalcular a assinatura | Não tem o `SECRET` — impossível |
| Reusar a assinatura antiga com payload novo | Assinatura é de **outro** payload — não confere |

> Qualquer alteração de 1 byte no header ou payload **invalida** a assinatura. É isso que torna o JWT confiável mesmo sendo legível: você pode **ler**, mas não pode **forjar**.

---

# HS256 vs RS256

| | HS256 (simétrico) | RS256 (assimétrico) |
|---|---|---|
| Chaves | **Um** segredo compartilhado | Par: privada assina, pública verifica |
| Quem pode emitir | Quem tem o segredo | Só quem tem a chave **privada** |
| Quem pode verificar | Quem tem o segredo | Qualquer um com a chave **pública** |
| Quando usar | API monolítica, um emissor | Vários serviços verificam, um emite |

```go
// HS256 — o que usamos na Sprint 3
token.SignedString([]byte(secret))
```

> Para o projeto, **HS256 basta**: vocês são emissor e verificador ao mesmo tempo. RS256 entra quando o emissor (ex.: um Identity Provider) é separado de quem consome o token.

---

# O ataque `alg: none`

Um JWT mal-validado aceita um header assim:

```json
{ "alg": "none", "typ": "JWT" }
```

`alg: none` significa "token sem assinatura". Uma biblioteca ingênua **aceita** — e aí qualquer um forja qualquer payload.

```go
// ✅ sempre declare o método esperado e cheque-o
token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("alg inesperado: %v", t.Header["alg"])
    }
    return []byte(secret), nil
})
```

> Regra de ouro: **o servidor decide qual algoritmo aceita**, nunca o token. Confiar no `alg` do header é confiar no atacante.

---

# O fluxo completo

```
1. POST /login  {email, senha}
        │
        ▼
2. Servidor: confere senha (bcrypt.CompareHashAndPassword)
        │
        ▼
3. Servidor: gera JWT assinado  →  devolve no corpo da resposta
        │
        ▼
4. Cliente guarda o token e o envia em CADA requisição:
        Authorization: Bearer eyJhbGc...
        │
        ▼
5. Middleware: valida assinatura + exp  →  extrai claims
        │
        ▼
6. Handler: já sabe quem é o usuário (via context)
```

---

# Passo 2-3: o handler de login

```go
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    json.NewDecoder(r.Body).Decode(&req)

    user, err := h.repo.FindByEmail(r.Context(), req.Email)
    if err != nil {
        http.Error(w, "credenciais inválidas", http.StatusUnauthorized)
        return
    }
    // senha guardada como HASH — nunca em texto puro
    err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password))
    if err != nil {
        http.Error(w, "credenciais inválidas", http.StatusUnauthorized)
        return
    }
    token, _ := generateToken(user)
    writeJSON(w, map[string]string{"access_token": token})
}
```

> Mensagem de erro **genérica** ("credenciais inválidas") tanto para email inexistente quanto para senha errada. Diferenciar as duas vaza informação: revela quais emails existem.

---

# Gerando o token

```go
func generateToken(user User) (string, error) {
    claims := jwt.MapClaims{
        "sub":  fmt.Sprint(user.ID),
        "name": user.Name,
        "role": user.Role,
        "iat":  time.Now().Unix(),
        "exp":  time.Now().Add(15 * time.Minute).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}
```

Pontos a observar:

- `exp` com **15 minutos** — token de acesso é curto por design
- `jwtSecret` vem de **variável de ambiente**, nunca hard-coded no repositório
- `role` na claim — é o que o middleware de autorização vai ler depois

---

# Passo 5: o middleware de autenticação

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")
        tokenStr, ok := strings.CutPrefix(header, "Bearer ")
        if !ok {
            http.Error(w, "token ausente", http.StatusUnauthorized)
            return
        }
        claims, err := validateToken(tokenStr)
        if err != nil {
            http.Error(w, "token inválido", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), userKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))   // ← segue o fluxo
    })
}
```

> O middleware é o **portão**. Se o token não é válido, a requisição **morre aqui** — o handler protegido nunca executa. Repare no padrão `func(http.Handler) http.Handler` da Sprint 0.

---

# Validando o token

```go
func validateToken(tokenStr string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("alg inesperado")
        }
        return []byte(jwtSecret), nil
    })
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("token inválido")
    }
    return token.Claims.(jwt.MapClaims), nil
}
```

A chamada `jwt.Parse` valida, em uma só passada: a **assinatura**, o **`exp`** (token expirado → erro) e o **`alg`** (via callback). Se qualquer uma falha, `token.Valid` é `false`.

---

# Aplicando o middleware com Chi

```go
r := chi.NewRouter()

// rotas públicas
r.Post("/login", authHandler.Login)
r.Post("/register", authHandler.Register)

// rotas protegidas — grupo com o middleware
r.Group(func(r chi.Router) {
    r.Use(AuthMiddleware)            // ← portão para tudo abaixo
    r.Get("/contacts", listContacts)
    r.Post("/contacts", createContact)
    r.Delete("/contacts/{id}", deleteContact)
})
```

> Mesma técnica de `r.Group` da Sprint 1: um middleware aplicado a um **subconjunto** de rotas. Login e register **ficam de fora** — senão ninguém conseguiria autenticar para entrar.

---

# Lendo a identidade no handler

```go
func listContacts(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(userKey).(jwt.MapClaims)
    userID := claims["sub"].(string)
    role   := claims["role"].(string)

    // agora o handler decide:
    // - filtrar contatos só do userID?
    // - exigir role == "admin"?
    ...
}
```

| O middleware fez | O handler faz |
|---|---|
| Provou **quem** é (autenticação) | Decide **o que** pode (autorização) |
| Pôs as claims no `context` | Lê as claims e aplica a regra de negócio |

> Separação limpa: o middleware não conhece regra de negócio; o handler não reimplementa validação de token. Cada um faz uma coisa.

---

# Armadilha: o nome do parâmetro errado

`r.Context().Value(userKey)` só funciona se a **chave** for a mesma na escrita e na leitura.

```go
// ❌ string crua como chave — colide com outros pacotes
ctx := context.WithValue(r.Context(), "user", claims)

// ✅ tipo próprio, não exportado — chave única e segura
type contextKey string
const userKey contextKey = "user"
ctx := context.WithValue(r.Context(), userKey, claims)
```

> `context.WithValue` com chave `string` é um *code smell* que o `go vet` reclama. Use um tipo próprio: garante que nenhum outro pacote sobrescreva sua chave por acidente.

---

# Erros comuns

| Sintoma | Causa provável |
|---|---|
| Token "válido" mas servidor rejeita | Segredo diferente entre gerar e validar |
| Funciona local, falha no CI | `JWT_SECRET` não definido como env var no workflow |
| Qualquer token passa | Esqueceu de checar `token.Valid` ou o `alg` |
| `401` em rota que deveria ser pública | Middleware aplicado no grupo errado |
| Token nunca expira | Faltou a claim `exp` no `MapClaims` |
| `panic` ao ler claim | Type assertion sem o `, ok` — claim ausente |

> Antes de pedir ajuda: confira **o segredo** primeiro. Mais da metade dos `401` misteriosos são segredo divergente entre os dois lados.

---

# Senhas: bcrypt, nunca texto puro

```go
// no register — guardar
hash, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
// guarda `hash` no banco, descarta a senha

// no login — conferir
err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(senha))
```

- `bcrypt` é **lento de propósito** — dificulta ataques de força bruta
- Inclui um *salt* automático — duas senhas iguais geram hashes diferentes
- A senha em texto puro **nunca** toca o banco, nem o log

> Vazamento de banco acontece. Se as senhas estiverem em texto puro, é catástrofe. Em hash bcrypt, o atacante ainda tem muito trabalho pela frente. Isso conecta com o **OWASP** que veremos na aula 11.

---

# Checklist da aula

A Sprint 3 (auth) entrega <span class="tag">prazo sprint: 29/05</span>:

- [ ] `POST /login` confere senha com `bcrypt` e devolve JWT
- [ ] Token assinado com `HS256` e segredo em **env var**
- [ ] Claim `exp` presente — token de acesso curto
- [ ] `AuthMiddleware` valida assinatura, `exp` e `alg`
- [ ] Rotas sensíveis dentro de `r.Group` com o middleware
- [ ] Handler lê a identidade via `context`

> Se até o fim da aula o login devolve um token e uma rota protegida responde `401` sem ele e `200` com ele, **o principal da autenticação está pronto**.

---

# Referências

**JWT**

- [RFC 7519 — JSON Web Token](https://datatracker.ietf.org/doc/html/rfc7519) — a especificação
- [jwt.io](https://jwt.io/) — decodificar e inspecionar tokens
- [`golang-jwt/jwt`](https://github.com/golang-jwt/jwt) — a biblioteca usada no projeto

**Segurança**

- [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — hashing de senhas
- [OWASP — JWT Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html) — boas práticas
