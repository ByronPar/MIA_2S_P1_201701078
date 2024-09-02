package main

import (
	"backend/management"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Define your routes
	//r.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	r.HandleFunc("/comando", handlePost).Methods("POST") // RECIBIRA UNA CADENA DE TEXTO

	// Start the server
	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Lee el cuerpo de la solicitud utilizando io.ReadAll
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "No se pudo leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	// Convierte el cuerpo a una cadena de texto
	text := string(body)
	// Responde con la misma cadena de texto
	//fmt.Fprintf(w, "Texto recibido: %s", text)
	fmt.Fprintf(w, management.Analizar(text))
}
