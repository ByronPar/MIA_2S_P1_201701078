package main

import (
	"backend/functions"
	"backend/management"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	// Define tus rutas
	r.HandleFunc("/comando", handlePost).Methods("POST")

	// Define la ruta para el método GET
	r.HandleFunc("/archivos", handleArchivos).Methods("POST")

	r.HandleFunc("/file-content", getFileContentHandler).Methods("GET") // Cambiado a GET

	// Envolver el router con el middleware CORS
	wrappedHandler := corsMiddleware(r)

	// Inicia el servidor
	fmt.Println("Servidor escuchando en el puerto 8088")
	http.ListenAndServe(":8088", wrappedHandler) // Iniciar el servidor
}

func handleArchivos(w http.ResponseWriter, r *http.Request) {
	// Llamar a enableCors para habilitar CORS
	enableCors(w)

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
		jsonResult, err := functions.ListFilesAndDirs(text)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al listar archivos: %v", err), http.StatusInternalServerError)
			return
		}

		// Configurar el Content-Type como application/json
		w.Header().Set("Content-Type", "application/json")

		// Escribir el JSON generado en la respuesta
		w.Write([]byte(jsonResult))
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// enableCors establece los encabezados CORS en la respuesta
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")              // Permitir todos los orígenes
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Permitir métodos POST y OPTIONS
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Permitir encabezados específicos
}
func handlePost(w http.ResponseWriter, r *http.Request) {
	// Llamar a enableCors para habilitar CORS
	enableCors(w)

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
		fmt.Fprintf(w, management.Analizar(text))
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

// Middleware para manejar CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Permitir solicitudes de cualquier origen
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Métodos permitidos
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Encabezados permitidos

		// Manejar preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r) // Llamar al siguiente manejador
	})
}

// Manejador para obtener el contenido del archivo
func getFileContentHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
