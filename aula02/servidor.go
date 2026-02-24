package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func handlerPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func celsiusParaFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

func handlerCelsius(w http.ResponseWriter, r *http.Request) {
	valorStr := r.URL.Query().Get("valor")

	celsius := 0.0
	if valorStr != "" {
		celsius, _ = strconv.ParseFloat(valorStr, 64)
	}

	fmt.Fprintf(w, "%.0f°F", celsiusParaFahrenheit(celsius))
}
