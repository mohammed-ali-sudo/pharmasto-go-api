package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

		ok, err, token := services.SignIn(db, creds.Username, creds.Password)
		if err != nil || !ok {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		// 1. Set the Authorization header with the JWT.
		w.Header().Set("Authorization", "Bearer "+token)

		// 2. Set the Content-Type header.
		w.Header().Set("Content-Type", "application/json")

		// 3. Write the HTTP status code and send all set headers.
		w.WriteHeader(http.StatusOK)

		fmt.Println("", token)
		json.NewEncoder(w).Encode(map[string]string{"message": "signin successful"})
	}
}
