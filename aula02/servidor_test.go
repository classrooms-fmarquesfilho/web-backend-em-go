// servidor_test.go
package main

import (
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
}
