package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Function to Create Part
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

		part.ID = id

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(part)
	}
}

// GetPart Api call Handler
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
		parts, err := repository.ListParts()
		json.NewEncoder(w).Encode(parts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}
}

// UpdatePart Api Handler
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

// Patch Part api handler
func PatchPartHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		// Get the existing part to update it
		existingPart, err := repository.GetPart(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Apply updates to the existing part
		for key, value := range updates {
			switch key {
			case "name":
				existingPart.Name = value.(string)
			case "price":
				existingPart.Price = value.(float64)
			case "description":
				existingPart.Description = value.(string)
			case "attributes":
				existingPart.Attributes = value.(map[string]string)
			case "images":
				existingPart.Images = value.([]string)
			case "sku":
				existingPart.SKU = value.(string)
			case "fitment_data":
				existingPart.FitmentData = value.([]string)
			case "location":
				existingPart.Location = value.(string)
			case "shipment":
				existingPart.Shipment = value.(ShipmentInfo)
			case "metadata":
				existingPart.Metadata = value.(map[string]string)
			}
		}

		if err := repository.UpdatePart(id, existingPart); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// DeletePart Handler
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

// Get Version Handler
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

// List Part version Handler
func ListPartVersionsHandler(repository *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		versions, err := repository.ListPartVersions(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(versions)
	}
}
