package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// SendJSONResponse transforms a value of any type into JSON and sends the
// JSON data as an HTTP response.
// If the value v cannot be properly coded into JSON the function returns an error.
func SendJSONResponse(w http.ResponseWriter, statusCode int, v any) error {
	// Convert the response value to JSON.
	jsonResponse, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("could not marshal value into JSON: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResponse)
	if err != nil {
		return fmt.Errorf("could not write to http.ResponseWriter: %w", err)
	}
	return nil
}

// HandleServerError sends an error message and stack trace to the error logger and
// then sends a generic 500 Internal Server Error response to the client.
func HandleServerError(w http.ResponseWriter, err error, errorLogger *log.Logger) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	_ = errorLogger.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}