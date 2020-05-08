package web

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONHeader is the header for JSON responses
const JSONHeader = "application/json; charset=UTF-8"

// StandardResponse is the JSON response from an API method, indicating success or failure.
type StandardResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RecordResponse the JSON response from a get call
type RecordResponse struct {
	StandardResponse
	Record interface{} `json:"record"`
}

// RecordsResponse the JSON response from a list call
type RecordsResponse struct {
	StandardResponse
	Records interface{} `json:"records"`
}

// formatStandardResponse returns a JSON response from an API method, indicating success or failure
func formatStandardResponse(code, message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", JSONHeader)
	response := StandardResponse{Code: code, Message: message}

	if len(code) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Encode the response as JSON
	encodeResponse(w, response)
}

// formatRecordResponse returns a JSON response from an api call
func formatRecordResponse(record interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", JSONHeader)
	response := RecordResponse{StandardResponse{}, record}

	// Encode the response as JSON
	encodeResponse(w, response)
}

// formatRecordsResponse returns a JSON response from an api call
func formatRecordsResponse(records interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", JSONHeader)
	response := RecordsResponse{StandardResponse{}, records}

	// Encode the response as JSON
	encodeResponse(w, response)
}

func encodeResponse(w http.ResponseWriter, response interface{}) {
	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error forming the response:", err)
	}
}
