package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/lista04/ex01/handler"
)

// ── Setup ───────────────────────────────────────────────────────────────────
//
// Os testes usam o PostgreSQL real do CI (ou o local do Codespace).
// Antes de cada teste, TRUNCATEia a tabela para garantir isolamento.

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Fatal("DATABASE_URL não está definida — rode com docker compose up ou no Codespace")
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("conectar ao Postgres: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func resetDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"TRUNCATE contacts RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func newServer(t *testing.T) http.Handler {
	t.Helper()
	pool := newPool(t)
	resetDB(t, pool)
	return handler.NewRouterWithPool(pool)
}

// ── Casos ───────────────────────────────────────────────────────────────────

func TestListContactsEmpty(t *testing.T) {
	srv := newServer(t)
	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /contacts vazio: esperava 200, recebeu %d (body=%s)", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type: esperava application/json, recebeu %q", ct)
	}
	body := strings.TrimSpace(rec.Body.String())
	// O JSON deve ser [] (não null nem objeto).
	if body != "[]" {
		t.Errorf("Esperava body \"[]\", recebeu %q", body)
	}
}

func TestCreateContact(t *testing.T) {
	srv := newServer(t)

	body := `{"name":"Maria","email":"maria@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /contacts: esperava 201, recebeu %d (body=%s)", rec.Code, rec.Body.String())
	}

	var created map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if _, ok := created["id"]; !ok {
		t.Errorf("Resposta sem campo \"id\": %v", created)
	}
	if created["name"] != "Maria" {
		t.Errorf("name esperado \"Maria\", recebeu %v", created["name"])
	}
	if created["email"] != "maria@example.com" {
		t.Errorf("email esperado, recebeu %v", created["email"])
	}
	if _, ok := created["created_at"]; !ok {
		t.Errorf("Resposta sem campo \"created_at\": %v", created)
	}
}

func TestCreateInvalidJSON(t *testing.T) {
	srv := newServer(t)
	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString("não é json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("JSON inválido: esperava 400, recebeu %d", rec.Code)
	}
}

func TestCreateMissingFields(t *testing.T) {
	srv := newServer(t)
	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(`{"name":"sem email"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Campos faltando: esperava 422, recebeu %d", rec.Code)
	}
}

func TestGetContact(t *testing.T) {
	srv := newServer(t)

	// Criar
	req := httptest.NewRequest(http.MethodPost, "/contacts",
		bytes.NewBufferString(`{"name":"João","email":"joao@x.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("setup falhou: %d", rec.Code)
	}
	var created map[string]any
	json.NewDecoder(rec.Body).Decode(&created)
	idFloat, ok := created["id"].(float64)
	if !ok {
		t.Fatalf("id não é número: %v", created["id"])
	}

	// Buscar
	url := "/contacts/" + intToStr(int(idFloat))
	req = httptest.NewRequest(http.MethodGet, url, nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s: esperava 200, recebeu %d (body=%s)", url, rec.Code, rec.Body.String())
	}
	var got map[string]any
	json.NewDecoder(rec.Body).Decode(&got)
	if got["name"] != "João" {
		t.Errorf("name retornado: %v", got["name"])
	}
}

func TestGetNotFound(t *testing.T) {
	srv := newServer(t)
	req := httptest.NewRequest(http.MethodGet, "/contacts/9999", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET inexistente: esperava 404, recebeu %d", rec.Code)
	}
}

func TestGetInvalidID(t *testing.T) {
	srv := newServer(t)
	req := httptest.NewRequest(http.MethodGet, "/contacts/abc", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("GET id inválido: esperava 400, recebeu %d", rec.Code)
	}
}

func TestDeleteContact(t *testing.T) {
	srv := newServer(t)

	// Criar
	req := httptest.NewRequest(http.MethodPost, "/contacts",
		bytes.NewBufferString(`{"name":"Ana","email":"ana@x.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	var created map[string]any
	json.NewDecoder(rec.Body).Decode(&created)
	idFloat := created["id"].(float64)
	url := "/contacts/" + intToStr(int(idFloat))

	// Deletar
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("DELETE: esperava 204, recebeu %d", rec.Code)
	}

	// Buscar deve dar 404
	req = httptest.NewRequest(http.MethodGet, url, nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET após DELETE: esperava 404, recebeu %d", rec.Code)
	}
}

func TestDeleteNotFound(t *testing.T) {
	srv := newServer(t)
	req := httptest.NewRequest(http.MethodDelete, "/contacts/9999", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("DELETE inexistente: esperava 404, recebeu %d", rec.Code)
	}
}

func TestListAfterCreating(t *testing.T) {
	srv := newServer(t)

	for _, b := range []string{
		`{"name":"A","email":"a@x.com"}`,
		`{"name":"B","email":"b@x.com"}`,
		`{"name":"C","email":"c@x.com"}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("setup: %d (%s)", rec.Code, rec.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET: %d", rec.Code)
	}
	var list []map[string]any
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 3 {
		t.Errorf("esperava 3 itens, recebeu %d", len(list))
	}
}

// ── helper local ────────────────────────────────────────────────────────────

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	s := string(digits)
	if neg {
		s = "-" + s
	}
	return s
}
