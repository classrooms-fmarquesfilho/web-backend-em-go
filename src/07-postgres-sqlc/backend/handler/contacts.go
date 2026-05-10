// Exercício 1 — Migrar CRUD de Contatos para sqlc + PostgreSQL
//
// SOLUÇÃO. Veja SOLUCAO.md (na pasta do gabarito) para o raciocínio
// passo-a-passo de cada um dos 4 TODOs.
//
// ── Endpoints ────────────────────────────────────────────────────────────────
//   GET    /contacts        → 200 + array JSON
//   POST   /contacts        → 201 + objeto criado (com id e created_at)
//   GET    /contacts/{id}   → 200 ou 404 (id é inteiro)
//   DELETE /contacts/{id}   → 204 ou 404

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/lista04/ex01/internal/db"
)

// ── Helpers de resposta (RFC 7807 + JSON) ───────────────────────────────────

type problemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func writeProblem(w http.ResponseWriter, status int, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problemDetails{
		Type:   "https://webii.ufrn.br/errors/" + http.StatusText(status),
		Title:  title,
		Status: status,
		Detail: detail,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// ── Aplicação ───────────────────────────────────────────────────────────────

// App agrupa o router e o cliente sqlc.
//
// Comparado à Sprint 1, sumiu o `map[int]Contact` e o mutex.
// No lugar entrou *db.Queries — todo o estado vive no PostgreSQL.
type App struct {
	queries *db.Queries
}

// NewRouter conecta ao Postgres e devolve o router já configurado.
//
// Usa DATABASE_URL do ambiente. Em produção, você nunca colocaria string
// de conexão hard-coded — sempre via variável de ambiente ou secret manager.
func NewRouter() http.Handler {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("não foi possível conectar ao banco: " + err.Error())
	}
	return NewRouterWithPool(pool)
}

// NewRouterWithPool é a variante usada pelos testes — recebe o pool já pronto.
func NewRouterWithPool(pool *pgxpool.Pool) http.Handler {
	a := &App{queries: db.New(pool)}

	r := chi.NewRouter()
	r.Get("/contacts", a.listContacts)
	r.Post("/contacts", a.createContact)
	r.Get("/contacts/{id}", a.getContact)
	r.Delete("/contacts/{id}", a.deleteContact)
	return r
}

// ── Handlers ────────────────────────────────────────────────────────────────

// listContacts retorna todos os contatos ordenados por created_at DESC.
func (a *App) listContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := a.queries.ListContacts(r.Context())
	if err != nil {
		writeProblem(w, http.StatusInternalServerError,
			"Internal Server Error", "failed to list contacts")
		return
	}

	// Quando não há linhas, ListContacts retorna `nil`. Sem essa
	// normalização, o JSON serializa como `null` em vez de `[]` — o
	// teste TestListContactsEmpty falha por causa disso.
	if contacts == nil {
		contacts = []db.Contact{}
	}

	writeJSON(w, http.StatusOK, contacts)
}

// createContact lê o JSON do body, valida campos obrigatórios e insere via sqlc.
func (a *App) createContact(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeProblem(w, http.StatusBadRequest, "Bad Request", "invalid json body")
		return
	}
	if body.Name == "" || body.Email == "" {
		writeProblem(w, http.StatusUnprocessableEntity,
			"Unprocessable Entity", "name and email are required")
		return
	}

	// CreateContact recebe um struct nomeado (db.CreateContactParams)
	// porque o INSERT tem mais de um parâmetro. Isso evita trocar
	// `name` com `email` por engano — um problema clássico de
	// parâmetros posicionais.
	contact, err := a.queries.CreateContact(r.Context(), db.CreateContactParams{
		Name:  body.Name,
		Email: body.Email,
	})
	if err != nil {
		writeProblem(w, http.StatusInternalServerError,
			"Internal Server Error", "failed to create contact")
		return
	}

	// Como a query usa `RETURNING *`, `contact` já vem com `id` e
	// `created_at` preenchidos pelo banco — é isso que vai no JSON.
	writeJSON(w, http.StatusCreated, contact)
}

// getContact busca um contato pelo ID (path param).
func (a *App) getContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeProblem(w, http.StatusBadRequest, "Bad Request", "id must be an integer")
		return
	}

	contact, err := a.queries.GetContact(r.Context(), int32(id))
	if err != nil {
		// pgx.ErrNoRows é o sentinel devolvido quando a query :one
		// não casa com nenhuma linha. Devemos diferenciar isso (404)
		// de outros erros de banco (500) — caso contrário o cliente
		// não sabe se o problema é dele (pediu id que não existe) ou
		// nosso (banco indisponível).
		if errors.Is(err, pgx.ErrNoRows) {
			writeProblem(w, http.StatusNotFound, "Not Found", "contact not found")
			return
		}
		writeProblem(w, http.StatusInternalServerError,
			"Internal Server Error", "failed to get contact")
		return
	}

	writeJSON(w, http.StatusOK, contact)
}

// deleteContact remove um contato pelo ID.
//
// A query `DELETE :exec` retorna apenas error — não diz se a linha
// existia. Para responder 404 corretamente, fazemos um GET prévio.
//
// (Há alternativas: trocar a anotação para :execrows, ou usar
// `DELETE ... RETURNING id` como :one. Para esta lista mantemos a
// versão simples e legível com duas queries.)
func (a *App) deleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeProblem(w, http.StatusBadRequest, "Bad Request", "id must be an integer")
		return
	}

	// 1. Verifica existência. Mesmo padrão do getContact.
	if _, err := a.queries.GetContact(r.Context(), int32(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeProblem(w, http.StatusNotFound, "Not Found", "contact not found")
			return
		}
		writeProblem(w, http.StatusInternalServerError,
			"Internal Server Error", "failed to check contact")
		return
	}

	// 2. Existe — pode deletar.
	if err := a.queries.DeleteContact(r.Context(), int32(id)); err != nil {
		writeProblem(w, http.StatusInternalServerError,
			"Internal Server Error", "failed to delete contact")
		return
	}

	// 3. 204 No Content: convenção REST para DELETE bem-sucedido sem
	//    corpo de resposta. Não use writeJSON aqui — 204 não pode ter
	//    body por especificação.
	w.WriteHeader(http.StatusNoContent)
}
