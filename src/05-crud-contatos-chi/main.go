// Exercício 1 — Migrar CRUD de Contatos para Chi
//
// Reimplemente o CRUD de contatos da Lista 2 (ex01) usando o router Chi
// no lugar do http.ServeMux. A interface pública e o comportamento devem
// ser EXATAMENTE OS MESMOS dos testes.
//
// Endpoints esperados:
//   GET    /contacts        → 200 + array JSON
//   POST   /contacts        → 201 + objeto criado
//   GET    /contacts/{id}   → 200 ou 404
//   DELETE /contacts/{id}   → 204 ou 404
//
// Requisitos:
//   - Usar github.com/go-chi/chi/v5 (já no go.mod)
//   - NewRouter() retorna http.Handler
//   - chi.URLParam(r, "id") para extrair parâmetros de path
//   - Content-Type "application/json" em todas as respostas com corpo
//
// NÃO chame http.ListenAndServe — os testes usam httptest.

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Contact é o modelo do recurso.
type Contact struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// app encapsula o estado da aplicação. NÃO use variáveis globais
// para storage — os testes criam uma nova instância para cada caso.
type app struct {
	contacts map[string]Contact
	nextID   int
}

// NewRouter retorna o router Chi configurado com todas as rotas.
// É a função pública chamada pelos testes.
func NewRouter() http.Handler {
	a := &app{
		contacts: make(map[string]Contact),
		nextID:   1,
	}

	r := chi.NewRouter()

	// TODO: registre as 4 rotas aqui usando r.Get, r.Post, r.Delete

	_ = a // remova esta linha quando começar a usar 'a'
	return r
}

// TODO: implemente os handlers como métodos de *app:
//
// func (a *app) listContacts(w http.ResponseWriter, r *http.Request)
// func (a *app) createContact(w http.ResponseWriter, r *http.Request)
// func (a *app) getContact(w http.ResponseWriter, r *http.Request)
// func (a *app) deleteContact(w http.ResponseWriter, r *http.Request)
