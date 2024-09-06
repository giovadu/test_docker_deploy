package routes

import (
	"notifcations_server/http/handlers"

	"github.com/gorilla/mux"
)

var router *mux.Router

func Router() *mux.Router {
	router = mux.NewRouter()
	router.HandleFunc("/send-notification", handlers.SendNotificationHandler).Methods("POST")
	return router
}
