package main

import (
	"log"
	"net/http"

	"example.com/lista04/ex01/handler"
)

func main() {
	router := handler.NewRouter()
	log.Println("API rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
