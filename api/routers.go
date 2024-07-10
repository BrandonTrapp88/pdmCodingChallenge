package main

import (
	"github.com/gorilla/mux"
)

func NewRouter(repository *Repository) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/parts", CreatePartHandler(repository)).Methods("POST")
	router.HandleFunc("/parts/{id}", GetPartHandler(repository)).Methods("GET")
	router.HandleFunc("/parts", ListPartsHandler(repository)).Methods("GET")
	router.HandleFunc("/parts/{id}", UpdatePartHandler(repository)).Methods("PATCH") // Change PUT to PATCH
	router.HandleFunc("/parts/{id}", DeletePartHandler(repository)).Methods("DELETE")
	router.HandleFunc("/parts/{id}/version/{version}", GetPartVersionHandler(repository)).Methods("GET")
	router.HandleFunc("/parts/{id}/versions", ListPartVersionsHandler(repository)).Methods("GET")
	router.HandleFunc("/search", SearchPartsHandler(repository)).Methods("GET")

	return router
}
