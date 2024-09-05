package handlers

import (
	"net/http"
	"notifcations_server/services"
	"notifcations_server/utils/io/response"
)

func SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	deviceToken := r.FormValue("device_token") // Obtenido del cuerpo de la solicitud
	title := r.FormValue("title")
	body := r.FormValue("body")
	if deviceToken == "" || title == "" || body == "" {
		response.Error("device_token, title y body son requeridos", http.StatusBadRequest, w)
		return
	}
	// Llamada al servicio para enviar la notificación
	err := services.SendFirebaseNotification(deviceToken, title, body)
	if err != nil {
		response.Error(err.Error(), http.StatusInternalServerError, w)
		return
	}
	response.Success("Notificación enviada exitosamente", http.StatusOK, w)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Error("Method not allowed", http.StatusMethodNotAllowed, w)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Error("Route not found", http.StatusNotFound, w)
}
