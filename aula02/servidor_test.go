// servidor_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerPing(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	rec := httptest.NewRecorder()

	handlerPing(rec, req)

	resultado := rec.Body.String()
	esperado := "pong"

	if resultado != esperado {
		t.Errorf("resultado %q, esperado %q", resultado, esperado)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("status: resultado %d, esperado %d", rec.Code, http.StatusOK)
	}
}

func TestHandlerCelsius(t *testing.T) {
	t.Run("ponto de congelamento", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/celsius?valor=0", nil)
		rec := httptest.NewRecorder()
		handlerCelsius(rec, req)
		verificarBody(t, rec, "32°F")
	})

	t.Run("ponto de ebulição", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/celsius?valor=100", nil)
		rec := httptest.NewRecorder()
		handlerCelsius(rec, req)
		verificarBody(t, rec, "212°F")
	})

	t.Run("sem parâmetro usa zero", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/celsius", nil)
		rec := httptest.NewRecorder()
		handlerCelsius(rec, req)
		verificarBody(t, rec, "32°F")
	})
}

func verificarBody(t testing.TB, rec *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	resultado := rec.Body.String()
	if resultado != esperado {
		t.Errorf("body: resultado %q, esperado %q", resultado, esperado)
	}
}
