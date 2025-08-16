package handlers

import (
	"database/sql"
	"encoding/json"
	"goapi/internal/services"
	"goapi/models"
	"net/http"
)

func InventoryCreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var m models.Manfactory
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&m); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if msg, ok := models.CreateManfactoryValidator(m); !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": "error",
				"error":  msg, // single error message
			})
			return
		}

		created, err := services.CreateManfactory(db, m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	}
}
