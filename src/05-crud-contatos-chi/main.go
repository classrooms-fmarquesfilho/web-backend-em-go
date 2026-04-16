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
	"encoding/json"
	"fmt"
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
	r.Get("/contacts", a.listContacts)
	r.Post("/contacts", a.createContact)
	r.Get("/contacts/{id}", a.getContact)
	r.Delete("/contacts/{id}", a.deleteContact)
	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (a *app) listContacts(w http.ResponseWriter, r *http.Request) {
	list := make([]Contact, 0, len(a.contacts))
	for _, c := range a.contacts {
		list = append(list, c)
	}
	writeJSON(w, http.StatusOK, list)
}

func (a *app) createContact(w http.ResponseWriter, r *http.Request) {
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	c.ID = fmt.Sprintf("%d", a.nextID)
	a.nextID++
	a.contacts[c.ID] = c
	writeJSON(w, http.StatusCreated, c)
}

func (a *app) getContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, ok := a.contacts[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (a *app) deleteContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := a.contacts[id]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	delete(a.contacts, id)
	w.WriteHeader(http.StatusNoContent)
}
