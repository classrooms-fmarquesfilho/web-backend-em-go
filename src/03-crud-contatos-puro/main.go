package main

// Exercício 1: CRUD de Contatos com Go 1.22+ ServeMux (25 pontos)
//
// Implemente os 4 handlers abaixo. Use os novos padrões de roteamento do Go 1.22+
// (método no pattern e r.PathValue para parâmetros de rota).
//
// Endpoints esperados:
//   GET    /contacts      → Lista todos os contatos (JSON array, status 200)
//   POST   /contacts      → Cria um contato (recebe JSON, status 201)
//   GET    /contacts/{id} → Busca contato por ID (status 200 ou 404)
//   DELETE /contacts/{id} → Remove contato por ID (status 204 ou 404)
//
// Regras:
//   - Armazene os contatos em memória (use uma map[string]Contact)
//   - IDs são strings sequenciais ("1", "2", "3", ...)
//   - GET /contacts retorna [] (array vazio) quando não há contatos, NÃO null
//   - Content-Type deve ser "application/json" em todas as respostas JSON
//
// IMPORTANTE: Os testes chamam NewRouter() diretamente. NÃO use http.DefaultServeMux
// e NÃO chame http.ListenAndServe.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Contact representa um contato armazenado pela API.
// Use este tipo no seu armazenamento e nas respostas JSON.
type Contact struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	var (
		mu     sync.Mutex
		store  = make(map[string]Contact)
		nextID int
	)

	mux.HandleFunc("GET /contacts", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		contacts := make([]Contact, 0, len(store))
		for _, c := range store {
			contacts = append(contacts, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contacts)
	})

	mux.HandleFunc("POST /contacts", func(w http.ResponseWriter, r *http.Request) {
		var c Contact
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		mu.Lock()
		nextID++
		c.ID = fmt.Sprintf("%d", nextID)
		store[c.ID] = c
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	})

	mux.HandleFunc("GET /contacts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		mu.Lock()
		c, ok := store[id]
		mu.Unlock()

		if !ok {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(c)
	})

	mux.HandleFunc("DELETE /contacts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		mu.Lock()
		_, ok := store[id]
		if ok {
			delete(store, id)
		}
		mu.Unlock()

		if !ok {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}
