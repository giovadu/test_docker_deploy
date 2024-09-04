package common

import (
	"net/http"
)

func SendReponse(writer http.ResponseWriter, status int, data []byte) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(data)
}

func SendError(writer http.ResponseWriter, status int) {
	data := []byte(`{"error": "Internal Server Error"}`)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(data)
}
