package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mustRouter(t *testing.T) *http.ServeMux {
	t.Helper()
	r := NewRouter()
	if r == nil {
		t.Fatal("NewRouter() retornou nil — você precisa implementar a função")
	}
	return r
}

func TestListContactsEmpty(t *testing.T) {
	router := mustRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /contacts: esperava 200, recebeu %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: esperava application/json, recebeu %q", ct)
	}
	var contacts []Contact
	if err := json.NewDecoder(rec.Body).Decode(&contacts); err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}
	if len(contacts) != 0 {
		t.Errorf("Esperava array vazio, recebeu %d contatos", len(contacts))
	}
}

func TestCreateAndGetContact(t *testing.T) {
	router := mustRouter(t)

	body := `{"name":"Maria","email":"maria@email.com"}`
	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /contacts: esperava 201, recebeu %d", rec.Code)
	}
	var created Contact
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}
	if created.ID == "" {
		t.Fatal("Contato criado sem ID")
	}
	if created.Name != "Maria" {
		t.Errorf("Nome: esperava Maria, recebeu %q", created.Name)
	}

	req = httptest.NewRequest(http.MethodGet, "/contacts/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /contacts/%s: esperava 200, recebeu %d", created.ID, rec.Code)
	}
	var fetched Contact
	json.NewDecoder(rec.Body).Decode(&fetched)
	if fetched.Email != "maria@email.com" {
		t.Errorf("Email: esperava maria@email.com, recebeu %q", fetched.Email)
	}
}

func TestGetContactNotFound(t *testing.T) {
	router := mustRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/contacts/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("GET /contacts/999: esperava 404, recebeu %d", rec.Code)
	}
}

func TestDeleteContact(t *testing.T) {
	router := mustRouter(t)

	body := `{"name":"João","email":"joao@email.com"}`
	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var created Contact
	json.NewDecoder(rec.Body).Decode(&created)

	if created.ID == "" {
		t.Fatal("POST /contacts não retornou ID — implemente POST primeiro")
	}

	req = httptest.NewRequest(http.MethodDelete, "/contacts/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("DELETE: esperava 204, recebeu %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/contacts/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET após DELETE: esperava 404, recebeu %d", rec.Code)
	}
}

func TestDeleteNotFound(t *testing.T) {
	router := mustRouter(t)
	req := httptest.NewRequest(http.MethodDelete, "/contacts/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("DELETE /contacts/999: esperava 404, recebeu %d", rec.Code)
	}
}

func TestListAfterCreating(t *testing.T) {
	router := mustRouter(t)

	for _, name := range []string{"Ana", "Carlos"} {
		b, _ := json.Marshal(map[string]string{"name": name, "email": name + "@test.com"})
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var contacts []Contact
	body, _ := io.ReadAll(rec.Body)
	json.Unmarshal(body, &contacts)
	if len(contacts) != 2 {
		t.Errorf("Esperava 2 contatos, recebeu %d", len(contacts))
	}
}
