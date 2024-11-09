package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

// errorMessageBody the standard data type used to deliver error messages.
type errorMessageBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	} `json:"error"`
}

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

	sendError(w, http.StatusInternalServerError)
}

func HandleNotFoundError(w http.ResponseWriter) {
	sendError(w, http.StatusNotFound)
}

func sendError(w http.ResponseWriter, statusCode int) {
	errorMessageBody := newErrorMessageBody(strconv.Itoa(statusCode), http.StatusText(statusCode), "")
	SendJSONResponse(w, statusCode, errorMessageBody)
}

// newErrorMessageBody returns a body for an error message that can be sent with
// SendJSONResponse.
func newErrorMessageBody(code, message, details string) errorMessageBody {
	body := errorMessageBody{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		}{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	return body
}

func SendErrorMessage(w http.ResponseWriter, statusCode int, errorLogger *log.Logger, errorMessage string) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := fmt.Fprintln(w, errorMessage)
	if err != nil {
		HandleServerError(w, fmt.Errorf("could not write to http.ResponseWriter: %w", err), errorLogger)
	}
}
