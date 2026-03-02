package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type requisicaoCalculo struct {
	A        float64 `json:"a"`
	B        float64 `json:"b"`
	Operacao string  `json:"operacao"`
}

func calcular(a float64, operacao string, b float64) float64 {
	switch operacao {
	case "soma":
		return a + b
	case "subtracao":
		return a - b
	case "multiplicacao":
		return a * b
	default:
		return 0
	}
}

func configurarRotas() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ajuda", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Endpoints disponíveis: /ping, /celsius, /calcular"))
	})

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dados := map[string]interface{}{
			"online": true,
			"versao": "1.0",
		}
		json.NewEncoder(w).Encode(dados)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Não Encontrado"))
			return
		}
		w.Write([]byte("Calculadora API"))
	})

	return mux
}

func handlerCalcular(w http.ResponseWriter, r *http.Request) {
	var req requisicaoCalculo

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"erro":"JSON inválido"}`))
		return
	}

	resposta := map[string]float64{
		"resultado": calcular(req.A, req.Operacao, req.B),
	}
	json.NewEncoder(w).Encode(resposta)
}

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
