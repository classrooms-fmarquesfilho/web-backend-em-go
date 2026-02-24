// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", handlerPing)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
