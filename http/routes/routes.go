package routes

import (
	"net/http"
	"notifcations_server/http/handlers"

	"github.com/gorilla/mux"
)

var router *mux.Router

func Router() *mux.Router {
	router = mux.NewRouter()
	//handlers para errores 404 y 405
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	//endpoints para enviar notificaciones hacia dispositivos segregado en dos proyectos de firebase diferentes
	router.HandleFunc("/v1/send-notification", handlers.SendNotificationHandlerV1).Methods("POST")
	router.HandleFunc("/v2/send-notification", handlers.SendNotificationHandlerV2).Methods("POST")

	//este endpoint es el que usa gpsec que perfectamente puede ser migrado a /v1/send-notification
	router.HandleFunc("/send-notification", handlers.SendNotificationHandlerV1).Methods("POST")

	// este endpoint para enviar a varios dispositivos al mismo tiempo
	router.HandleFunc("/v1/send-bulk-notification", handlers.SendBulkNotificationHandlerV1).Methods("POST")
	router.HandleFunc("/v2/send-bulk-notification", handlers.SendBulkNotificationHandlerV2).Methods("POST")
	return router
}
