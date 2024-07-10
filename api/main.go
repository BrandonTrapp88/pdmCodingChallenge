package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
)

func main() {
	// Set up the database connection
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "vehicle_parts_db"

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize the repository with the database connection
	repository := NewRepository(db)
	router := NewRouter(repository)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE", "PATCH"})

	log.Println("Starting server on :1710")
	if err := http.ListenAndServe(":1710", handlers.CORS(originsOk, headersOk, methodsOk)(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
