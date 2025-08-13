package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"goapi/internal/services"
	"goapi/models"
	"goapi/shared"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

func Delethandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}
		err2 := services.Delet(db, id)
		fmt.Println("", err2)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"id":     id,
		})
	}
}

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

		if msg, ok := models.CreateTeacherValidator(t); !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": "error",
				"error":  msg, // single error message
			})
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

func TeacherGetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ✅ Extract path variable
		vars := mux.Vars(r)
		idStr := vars["id"]

		// ✅ Validate and convert ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		// ✅ Call service to get teacher
		teacher, err := services.GetTeacherByID(db, id)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// ✅ Return JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teacher)
	}
}

func TeachersGetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Define a struct to match JSON: { "ids": [...] }
		var req struct {
			IDs []int `json:"ids"`
		}

		// Decode the JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// Call the service
		teachers, err := services.GetTeachers(db, req.IDs)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teachers)
	}
}

func TeachersGetHandlerfilter(db *sql.DB) http.HandlerFunc {
	type response struct {
		Data []models.Teacher `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 1) Build filters from query params
		if len(r.URL.Query()) > 0 {
			filters := map[string]string{}
			for k, v := range r.URL.Query() {
				if len(v) > 0 {
					filters[k] = v[0]
				}
			}

			teachers, err := services.FilterServices(db, filters)
			if err != nil {
				http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			// Ensure non-nil slice
			if teachers == nil {
				teachers = []models.Teacher{}
			}
			json.NewEncoder(w).Encode(response{Data: teachers})
			return
		}

		// 2) No query → decode IDs from body
		var req struct {
			IDs []int `json:"ids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		teachers, err := services.GetTeachers(db, req.IDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if teachers == nil {
			teachers = []models.Teacher{}
		}
		json.NewEncoder(w).Encode(response{Data: teachers})
	}
}

func Filter(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if len(query) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Convert query to map
		filters := map[string]string{}
		for k, v := range query {
			if len(v) > 0 {
				filters[k] = v[0]
			}
		}

		teachers, err := services.FilterServices(db, filters)
		if err != nil {
			http.Error(w, "Filter error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return filtered result directly
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teachers)
	})
}

// Sort is middleware that checks for ?sortby=field[:asc|desc].
// If present, it runs services.SortServices and writes the JSON response.
// Otherwise it calls next.ServeHTTP.
func Sort(db *sql.DB, next http.Handler) http.Handler {
	type resp struct {
		Data []models.Teacher `json:"data"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sortParams := r.URL.Query()["sortby"]
		if len(sortParams) > 0 {
			// Call your SortServices
			teachers, err := services.SortServices(db, sortParams)
			if err != nil {
				// Bad sort param or DB error
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if teachers == nil {
				teachers = []models.Teacher{}
			}
			json.NewEncoder(w).Encode(resp{Data: teachers})
			return
		}

		// No sort param → pass through
		next.ServeHTTP(w, r)
	})
}

func TeacherUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow PUT
		if r.Method != http.MethodPut {
			http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1) Extract path variable
		vars := mux.Vars(r)
		idStr := vars["id"]

		// 2) Validate and convert ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		// 3) Decode request body into a Teacher struct
		var input models.Teacher
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		// 4) Call service to update and return the updated teacher
		updated, err := services.UpdateTeacherService(db, input, id)
		if err != nil {
			http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 5) Return JSON of the updated teacher
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	}
}

func TeacherPatchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Only PATCH method is allowed", http.StatusMethodNotAllowed)
			return
		}
		// Extract and validate ID
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		userData, err := services.GetTeacherByID(db, id)

		fmt.Println("user data", userData)
		// Decode incoming JSON into models.Teacher
		var inputTeacherFromBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&inputTeacherFromBody); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		teacherValue := reflect.ValueOf(&userData).Elem()
		fmt.Println("teacher value from userdata", teacherValue.Type().Field(0))
		fmt.Println("teacher value from userdata", teacherValue.Type().Field(1))

		for k, v := range inputTeacherFromBody {
			for i := 0; i < teacherValue.NumField(); i++ {

				fmt.Println("key", k)
				field := teacherValue.Type().Field(i)
				fmt.Println("", field.Tag.Get("json"))
				if field.Tag.Get("json") == k+" ,omitempty" {
					if teacherValue.Field(i).CanSet() {
						teacherValue.Field(i).Set(reflect.ValueOf(v).Convert(teacherValue.Field(i).Type()))
					}
				}
			}
		}

		// Call service to patch and get updated record
		updatedTeacherFromDb, err := services.PatchTeacherService(db, userData, id)
		if err != nil {
			http.Error(w, "Patch failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return updated teacher
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTeacherFromDb)
	}
}

// PatchTeachersHandler handles PATCH requests to update multiple teachers.
func PatchTeachersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Only PATCH method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON request body into a slice of maps.
		var teachers []map[string]any
		err := json.NewDecoder(r.Body).Decode(&teachers)
		if err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// Call the service function to perform the updates.
		err = services.PatchTeachersServices(db, teachers)
		if err != nil {
			// Handle different types of errors from the service.
			if err.Error() == "invalid id format or missing id" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err.Error() == "teacher not found" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			// For all other errors, return an internal server error.
			http.Error(w, fmt.Sprintf("failed to update teachers: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with a success message.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Teachers updated successfully",
		})
	}
}
