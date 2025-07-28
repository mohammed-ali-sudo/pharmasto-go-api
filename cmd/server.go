package main

import (
	"fmt"
	"goapi/internal/handlers"
	"goapi/shared/db"
	middleware "goapi/shared/middlewares"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Open DB
	connStr := "postgres://postgres:1010204080@database-postgressdb.clgkywaycm2o.eu-north-1.rds.amazonaws.com:5432/pharmasto?sslmode=require"
	database := db.Open(connStr)
	defer database.Close()

	// 2. Router + middleware
	router := mux.NewRouter()
	router.Use(middleware.CORS)
	router.Use(middleware.ResponseTimeMw)
	router.Use(middleware.SecurityHeader)

	// 3. Route using the simplified handler
	router.HandleFunc("/teacher", handlers.TeacherCreateHandler(database)).Methods("POST")
	router.HandleFunc("/teacher/all", handlers.TeachersGetHandler(database)).Methods("POST")
	router.HandleFunc("/teacher/{id}", handlers.TeacherGetHandler(database)).Methods("GET")

	// 4. Start server
	fmt.Println("🚀 Server listening on :8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}
