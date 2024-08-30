package routes

import (
	"net/http"

	"test_dcoker_deploy/http/handlers"

	"github.com/gorilla/mux"
)

var router *mux.Router

func Router() *mux.Router {
	router = mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	router.HandleFunc("/test", handlers.Test).Methods("GET")
	return router
}
