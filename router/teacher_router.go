package router

import (
	"database/sql"
	"goapi/internal/handlers"

	middleware "goapi/shared/middlewares"

	"github.com/gorilla/mux"
)

func TeacherRouter(database *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// --- Public routes (no SecurityHeader) ---
	public := router.PathPrefix("/").Subrouter()
	public.HandleFunc("/signup", handlers.SignUpHandler(database)).Methods("POST")
	public.HandleFunc("/signin", handlers.SignInHandler(database)).Methods("POST")

	// --- Private routes (with SecurityHeader) ---
	private := router.PathPrefix("/").Subrouter()
	private.Use(middleware.Protect)

	private.HandleFunc("/", handlers.TeacherCreateHandler(database)).Methods("POST")
	private.HandleFunc("/all", handlers.TeachersGetHandler(database)).Methods("POST")
	private.HandleFunc("/all", handlers.PatchTeachersHandler(database)).Methods("PATCH")
	private.HandleFunc("/filter", handlers.TeachersGetHandlerfilter(database)).Methods("POST", "GET")
	private.HandleFunc("/{id}", handlers.Delethandler(database)).Methods("DELETE")
	private.HandleFunc("/{id}", handlers.TeacherGetHandler(database)).Methods("GET")
	private.HandleFunc("/{id}", handlers.TeacherUpdateHandler(database)).Methods("PUT")
	private.HandleFunc("/{id}", handlers.TeacherPatchHandler(database)).Methods("PATCH")
	private.HandleFunc("/manfactory", handlers.InventoryCreateHandler(database)).Methods("POST")

	return router
}
