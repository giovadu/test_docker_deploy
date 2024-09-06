package routes

import (
	"notifcations_server/http/handlers"

	"github.com/gorilla/mux"
)

var router *mux.Router

func Router() *mux.Router {
	router = mux.NewRouter()
	router.HandleFunc("/v1/send-notification", handlers.SendNotificationHandlerV1).Methods("POST")
	router.HandleFunc("/v2/send-notification", handlers.SendNotificationHandlerV2).Methods("POST")
	router.HandleFunc("/send-notification", handlers.SendNotificationHandlerV1).Methods("POST")
	return router
}
