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
	fmt.Println("Server is listening on port 8088")
	log.Fatal(http.ListenAndServe(":8088", r))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")              // Permitir todos los orígenes
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Permitir métodos POST y OPTIONS
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Permitir encabezados específicos
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	enableCors(&w) // Habilitar CORS para todas las respuestas
	// Lee el cuerpo de la solicitud utilizando io.ReadAll

	if r.Method == http.MethodOptions {
		// Manejar las solicitudes preflight
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		// Aquí procesas el comando
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
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}
