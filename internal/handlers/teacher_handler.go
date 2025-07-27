package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"goapi/internal/services"
	"goapi/models"
)

func TeacherCreateHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
            return
        }

        var t models.Teacher
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
            return
        }

        created, err := services.AddTeacher(db, t)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(created)
    }
}