package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"goapi/models"
	sqlconnect "goapi/repo"
)

type Teacher = models.Teacher

var (
	mutex    = &sync.Mutex{}
	teachers = make(map[int]Teacher)
)

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Class     string `json:"class"`
		Subject   string `json:"subject"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.Connectdb("gotest")
	if err != nil {
		http.Error(w, "❌ Failed to connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "❌ Failed to prepare SQL", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	email := fmt.Sprintf("%s.%s@example.com", input.FirstName, input.LastName)

	res, err := stmt.Exec(input.FirstName, input.LastName, email, input.Class, input.Subject)
	if err != nil {
		http.Error(w, "❌ Failed to execute insert", http.StatusInternalServerError)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "❌ Failed to get last insert ID", http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	t := Teacher{
		ID:        int(lastID),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     email,
		Class:     input.Class,
		Subject:   input.Subject,
	}
	teachers[t.ID] = t

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status  string  `json:"status"`
		Teacher Teacher `json:"teacher"`
	}{
		Status:  "success",
		Teacher: t,
	})
}

func main() {
	db, err := sqlconnect.Connectdb("gotest")
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	if !testDBConnection(db) {
		log.Fatal("❌ Database connection test failed.")
	}
	fmt.Println("✅ Database connection verified.")

	mux := http.NewServeMux()
	mux.HandleFunc("/teacher", addTeacherHandler)

	server := &http.Server{
		Addr:    ":8001",
		Handler: mux,
	}

	fmt.Println("🚀 Server is running on http://localhost:8001")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}

func testDBConnection(db *sql.DB) bool {
	row := db.QueryRow("SELECT 1")
	var dummy int
	err := row.Scan(&dummy)
	if err != nil || dummy != 1 {
		fmt.Println("❌ Error pinging DB or invalid result:", err)
		return false
	}
	return true
}

type Middleware func(http.Handler) http.Handler

func ApplyMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}
