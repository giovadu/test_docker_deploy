package handlers

import (
	"log"
	"net/http"
	"test_dcoker_deploy/utils/io/response"
)

func Test(w http.ResponseWriter, r *http.Request) {
	log.Printf("saludos desde la ruta test")
	response.Success(nil, http.StatusOK, w)
}
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Error("Method not allowed", http.StatusMethodNotAllowed, w)
}
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Error("Route not found", http.StatusNotFound, w)
}
