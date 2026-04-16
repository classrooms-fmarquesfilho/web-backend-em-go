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
	"net/http"
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

	// Dica: você pode declarar o estado aqui (store map, contador de IDs, sync.Mutex)
	// var (
	//     mu     sync.Mutex
	//     store  = make(map[string]Contact)
	//     nextID int
	// )

	mux.HandleFunc("GET /contacts", func(w http.ResponseWriter, r *http.Request) {
		// IMPLEMENTAR: Retornar a lista de contatos como JSON array
		_ = Contact{} // remova esta linha quando implementar
	})

	mux.HandleFunc("POST /contacts", func(w http.ResponseWriter, r *http.Request) {
		// IMPLEMENTAR: Criar um novo contato (status 201)
	})

	mux.HandleFunc("GET /contacts/{id}", func(w http.ResponseWriter, r *http.Request) {
		// IMPLEMENTAR: Retornar um contato específico por ID
		// Use r.PathValue("id") para extrair o parâmetro da rota
	})

	mux.HandleFunc("DELETE /contacts/{id}", func(w http.ResponseWriter, r *http.Request) {
		// IMPLEMENTAR: Deletar um contato específico por ID (status 204)
	})

	return mux
}
