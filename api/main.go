package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	dsn := "root:blue1234@tcp(localhost:3306)/vehicle_parts_db"
	repository, err := NewRepository(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := NewRouter(repository)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE", "Patch"})

	log.Println("Starting server on :1710")
	if err := http.ListenAndServe(":1710", handlers.CORS(originsOk, headersOk, methodsOk)(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
