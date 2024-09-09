package response

import (
	"encoding/json"
	"net/http"
)

// ResponseStruct struct to represent the error or done response format
type ResponseStruct struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func Error(message string, status int, w http.ResponseWriter) {
	// Create ErrorResponse struct
	errResponse := ResponseStruct{
		StatusCode: status,
		Message:    message,
		Data:       nil,
	}

	// Call sendResponse to send the JSON response
	sendResponse(errResponse, status, w)
}

func Success(data interface{}, status int, w http.ResponseWriter) {
	// Create SuccessResponse struct
	successResponse := ResponseStruct{
		Message:    "Success",
		StatusCode: status,
		Data:       data,
	}

	// Call sendResponse to send the JSON response
	sendResponse(successResponse, status, w)
}

func sendResponse(response interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
