package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

// SendError sends a JSON error response
func SendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status: "error",
		Msg:    message,
	})
}

func ErrorHandling(err error, msg string) error {

	errorLogger := log.New(os.Stderr, "error:", log.Ldate|log.Lshortfile)
	errorLogger.Println(msg, err)
	return fmt.Errorf(msg)
}
