package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePartHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var part Part
		if err := json.NewDecoder(r.Body).Decode(&part); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := repository.CreatePart(part)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		part.ID = id // Ensure the ID is set in the response

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(part)
	}
}

func GetPartHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		part, err := repository.GetPart(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(part)
	}
}

func ListPartsHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := repository.ListParts()
		json.NewEncoder(w).Encode(parts)
	}
}

func UpdatePartHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var part Part
		if err := json.NewDecoder(r.Body).Decode(&part); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repository.UpdatePart(id, part); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeletePartHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if err := repository.DeletePart(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetPartVersionHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		version, err := strconv.Atoi(mux.Vars(r)["version"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		part, err := repository.GetPartVersion(id, version)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(part)
	}
}

func SearchPartsHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		results := repository.SearchPartsByName(name)
		json.NewEncoder(w).Encode(results)
	}
}
