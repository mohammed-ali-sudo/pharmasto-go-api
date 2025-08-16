package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"goapi/internal/services"
	"goapi/models"
)

// POST /signup
func SignUpHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := services.SignUp(db, &user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
	}
}

// POST /signin
func SignInHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		ok, err := services.SignIn(db, creds.Username, creds.Password)
		if err != nil || !ok {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "signin successful"})
	}
}
