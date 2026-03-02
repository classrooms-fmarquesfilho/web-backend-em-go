// servidor_test.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRotas(t *testing.T) {
	mux := configurarRotas()

	t.Run("raiz", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		verificarBody(t, rec, "Calculadora API")
		verificarStatus(t, rec, http.StatusOK)
	})

	t.Run("ajuda", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ajuda", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		verificarBody(t, rec, "Endpoints disponíveis: /ping, /celsius, /calcular")
	})

	t.Run("status JSON", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/status", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		verificarStatus(t, rec, http.StatusOK)

		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type: resultado %q, esperado %q", contentType, "application/json")
		}
	})

	t.Run("não encontrado", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/xyz", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		verificarStatus(t, rec, http.StatusNotFound)
	})
}

func verificarStatus(t testing.TB, rec *httptest.ResponseRecorder, esperado int) {
	t.Helper()
	if rec.Code != esperado {
		t.Errorf("status: resultado %d, esperado %d", rec.Code, esperado)
	}
}

func TestHandlerCalcular(t *testing.T) {
	t.Run("soma válida", func(t *testing.T) {
		body := bytes.NewBufferString(`{"a":10,"operacao":"soma","b":5}`)
		req := httptest.NewRequest(http.MethodPost, "/calcular", body)
		rec := httptest.NewRecorder()

		handlerCalcular(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		verificarStatus(t, rec, http.StatusOK)

		var resultado map[string]interface{}
		respBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(respBody, &resultado)

		if resultado["resultado"] != 15.0 {
			t.Errorf("resultado %v, esperado 15", resultado["resultado"])
		}
	})

	t.Run("JSON inválido", func(t *testing.T) {
		body := bytes.NewBufferString("não é json")
		req := httptest.NewRequest(http.MethodPost, "/calcular", body)
		rec := httptest.NewRecorder()

		handlerCalcular(rec, req)

		verificarStatus(t, rec, http.StatusBadRequest)
	})
}

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
