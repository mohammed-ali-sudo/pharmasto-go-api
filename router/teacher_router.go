package router

import (
	"database/sql"
	"goapi/internal/handlers"

	"github.com/gorilla/mux"
)

func TeacherRouter(database *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Routes no longer need /teacher prefix
	router.HandleFunc("/", handlers.TeacherCreateHandler(database)).Methods("POST")
	router.HandleFunc("/all", handlers.TeachersGetHandler(database)).Methods("POST")
	router.HandleFunc("/all", handlers.PatchTeachersHandler(database)).Methods("PATCH")
	router.HandleFunc("/filter", handlers.TeachersGetHandlerfilter(database)).Methods("POST", "GET")
	router.HandleFunc("/{id}", handlers.Delethandler(database)).Methods("DELETE")
	router.HandleFunc("/{id}", handlers.TeacherGetHandler(database)).Methods("GET")
	router.HandleFunc("/{id}", handlers.TeacherUpdateHandler(database)).Methods("PUT")
	router.HandleFunc("/{id}", handlers.TeacherPatchHandler(database)).Methods("PATCH")
	router.HandleFunc("/manfactory", handlers.InventoryCreateHandler(database)).Methods("POST")
	router.HandleFunc("/signup", handlers.SignUpHandler(database)).Methods("POST")
	router.HandleFunc("/signin", handlers.SignInHandler(database)).Methods("POST")
	return router
}
