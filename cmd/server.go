package main

import (
	"fmt"
	"goapi/router" // package containing TeacherRouter
	"goapi/shared/db"
	middleware "goapi/shared/middlewares"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)



func main() {
	// Open DB
	connStr := "postgres://postgres:1010204080@database-postgressdb.clgkywaycm2o.eu-north-1.rds.amazonaws.com:5432/pharmasto?sslmode=require"
	database := db.Open(connStr)
	defer database.Close()

	// Router + middleware
	mainRouter := mux.NewRouter() // <- renamed from "router"
	mainRouter.Use(middleware.CORS)
	mainRouter.Use(middleware.ResponseTimeMw)
	mainRouter.Use(middleware.SecurityHeader)

	// Mount TeacherRouter under /teacher
	mainRouter.PathPrefix("/teacher").Handler(
		http.StripPrefix("/teacher", router.TeacherRouter(database)), // <- now this works
	)

	fmt.Println("🚀 Server listening on :8001")
	log.Fatal(http.ListenAndServe(":8001", mainRouter))

}
