package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Se NewRouter() ainda não foi implementado e retorna nil, o teste falha em vez
// de panic.
func newServer(t *testing.T) http.Handler {
	t.Helper()
	h := NewRouter()
	if h == nil {
		t.Fatal("NewRouter() retornou nil — implemente a função antes de rodar os testes")
	}
	return h
}

func TestListContactsEmpty(t *testing.T) {
	router := newServer(t)
	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /contacts: esperava 200, recebeu %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type: esperava application/json, recebeu %q", ct)
	}
	var contacts []Contact
	if err := json.NewDecoder(rec.Body).Decode(&contacts); err != nil {
		t.Fatalf("Erro ao decodificar JSON: %v", err)
	}
	if len(contacts) != 0 {
		t.Errorf("Esperava array vazio, recebeu %d contatos", len(contacts))
	}
}

func TestCreateAndGetContact(t *testing.T) {
	router := newServer(t)

	body := `{"name":"Maria","email":"maria@email.com"}`
	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /contacts: esperava 201, recebeu %d (body: %s)", rec.Code, rec.Body.String())
	}

	var created Contact
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("Erro ao decodificar resposta do POST: %v", err)
	}
	if created.ID == "" {
		t.Fatal("Esperava que o contato criado tivesse ID")
	}
	if created.Name != "Maria" || created.Email != "maria@email.com" {
		t.Errorf("Dados retornados não batem: %+v", created)
	}

	// Buscar por ID
	req = httptest.NewRequest(http.MethodGet, "/contacts/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /contacts/%s: esperava 200, recebeu %d", created.ID, rec.Code)
	}

	var fetched Contact
	if err := json.NewDecoder(rec.Body).Decode(&fetched); err != nil {
		t.Fatalf("Erro ao decodificar GET: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID retornado %q ≠ esperado %q", fetched.ID, created.ID)
	}
}

func TestGetContactNotFound(t *testing.T) {
	router := newServer(t)
	req := httptest.NewRequest(http.MethodGet, "/contacts/nao-existe", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET inexistente: esperava 404, recebeu %d", rec.Code)
	}
}

func TestDeleteContact(t *testing.T) {
	router := newServer(t)

	// Criar primeiro
	body := `{"name":"João","email":"joao@email.com"}`
	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var created Contact
	json.NewDecoder(rec.Body).Decode(&created)

	// Deletar
	req = httptest.NewRequest(http.MethodDelete, "/contacts/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("DELETE: esperava 204, recebeu %d", rec.Code)
	}

	// Confirmar que sumiu
	req = httptest.NewRequest(http.MethodGet, "/contacts/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET após delete: esperava 404, recebeu %d", rec.Code)
	}
}

func TestDeleteNotFound(t *testing.T) {
	router := newServer(t)
	req := httptest.NewRequest(http.MethodDelete, "/contacts/nao-existe", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("DELETE inexistente: esperava 404, recebeu %d", rec.Code)
	}
}

func TestListAfterCreating(t *testing.T) {
	router := newServer(t)

	bodies := []string{
		`{"name":"Ana","email":"ana@x.com"}`,
		`{"name":"Bruno","email":"bruno@x.com"}`,
	}
	for _, b := range bodies {
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("Setup falhou: %d", rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /contacts: esperava 200, recebeu %d", rec.Code)
	}
	var contacts []Contact
	if err := json.NewDecoder(rec.Body).Decode(&contacts); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(contacts) != 2 {
		t.Errorf("Esperava 2 contatos, recebeu %d", len(contacts))
	}
}
