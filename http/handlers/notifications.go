package handlers

import (
	"net/http"
	"notifcations_server/services"
	"notifcations_server/utils/response"
	"strings"
)

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Error("Method not allowed", http.StatusMethodNotAllowed, w)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Error("Route not found", http.StatusNotFound, w)
}

func SendNotificationHandlerV1(w http.ResponseWriter, r *http.Request) {
	deviceToken := r.FormValue("device_token") // Obtenido del cuerpo de la solicitud
	title := r.FormValue("title")
	body := r.FormValue("body")
	if deviceToken == "" || title == "" || body == "" {
		response.Error("device_token, title y body son requeridos", http.StatusBadRequest, w)
		return
	}
	// Llamada al servicio para enviar la notificación
	err := services.SendFirebaseNotificationV1(deviceToken, title, body)
	if err != nil {
		if strings.Contains(err.Error(), `"code": 404`) {
			response.Error("TOKEN UNREGISTERED", http.StatusNotFound, w)
			return
		}
		if strings.Contains(err.Error(), `"code": 400`) {
			response.Error("INVALID TOKEN", http.StatusNotFound, w)
			return
		}
		response.Error(err.Error(), http.StatusInternalServerError, w)
		return
	}
	response.Success("Notificación enviada exitosamente", http.StatusOK, w)
}
func SendNotificationHandlerV2(w http.ResponseWriter, r *http.Request) {
	deviceToken := r.FormValue("device_token") // Obtenido del cuerpo de la solicitud
	title := r.FormValue("title")
	body := r.FormValue("body")
	if deviceToken == "" || title == "" || body == "" {
		response.Error("device_token, title y body son requeridos", http.StatusBadRequest, w)
		return
	}
	// Llamada al servicio para enviar la notificación
	err := services.SendFirebaseNotificationV2(deviceToken, title, body)
	if err != nil {
		if strings.Contains(err.Error(), `"code": 404`) {
			response.Error("TOKEN UNREGISTERED", http.StatusNotFound, w)
			return
		}
		if strings.Contains(err.Error(), `"code": 400`) {
			response.Error("INVALID TOKEN", http.StatusNotFound, w)
			return
		}
		response.Error(err.Error(), http.StatusInternalServerError, w)
		return
	}
	response.Success("Notificación enviada exitosamente", http.StatusOK, w)
}
func SendBulkNotificationHandlerV1(w http.ResponseWriter, r *http.Request) {
	devicesTokens := r.FormValue("devices_tokens") // Obtenido del cuerpo de la solicitud
	title := r.FormValue("title")
	body := r.FormValue("body")
	if title == "" || body == "" {
		response.Error("title y body son requeridos", http.StatusBadRequest, w)
		return
	}
	if devicesTokens == "" {
		response.Error("el campo devices_tokens es requirdo", http.StatusBadRequest, w)
		return
	}
	// Llamada al servicio para enviar la notificación
	reponse, err := services.SendBulkFirebaseNotificationsV1(devicesTokens, title, body)
	if err != nil {
		response.Error(err.Error(), http.StatusInternalServerError, w)
		return
	}
	response.Success(reponse, http.StatusOK, w)
}
func SendBulkNotificationHandlerV2(w http.ResponseWriter, r *http.Request) {
	devicesTokens := r.FormValue("devices_tokens") // Obtenido del cuerpo de la solicitud
	title := r.FormValue("title")
	body := r.FormValue("body")
	if title == "" || body == "" {
		response.Error("title y body son requeridos", http.StatusBadRequest, w)
		return
	}
	if devicesTokens == "" {
		response.Error("el campo devices_tokens es requirdo", http.StatusBadRequest, w)
		return
	}
	// Llamada al servicio para enviar la notificación
	reponse, err := services.SendBulkFirebaseNotificationsV2(devicesTokens, title, body)
	if err != nil {
		response.Error(err.Error(), http.StatusInternalServerError, w)
		return
	}
	response.Success(reponse, http.StatusOK, w)
}
