package routes

import (
	"net/http"

	"notifcations_server/http/handlers"

	"github.com/gorilla/mux"
)

var router *mux.Router

func Router() *mux.Router {
	router = mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	router.HandleFunc("/send-notification", handlers.SendNotificationHandler).Methods("POST")
	return router
}
