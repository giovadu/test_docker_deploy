package services

import (
	"notifcations_server/database/repositories"
	"sync"
	"time"
)

// Global variables to store token and expiration
var (
	cachedToken     string
	tokenExpiryTime time.Time
	tokenLock       sync.Mutex
)

// GetAccessTokenWithCache obtiene el token de acceso, usando uno cacheado si es válido.
func GetAccessTokenWithCache() (string, error) {
	tokenLock.Lock()
	defer tokenLock.Unlock()

	// Verificar si el token ha expirado o no existe
	if time.Now().Before(tokenExpiryTime) && cachedToken != "" {
		return cachedToken, nil
	}

	// Obtener el TokenRefresher que maneja el refresco del token
	tokenRefresher, err := repositories.GetFirebaseClient("gd-notificacionesandroid-firebase-adminsdk-2v5rt-c75d589044.json")
	if err != nil {
		return "", err
	}

	// Obtener el nuevo token de acceso
	accessToken, err := tokenRefresher.GetAccessToken()
	if err != nil {
		return "", err
	}

	// Actualizar el token y su tiempo de expiración
	expireTime := time.Now().Add(time.Minute * 50) // Ajustar según la duración real del token
	cachedToken = accessToken
	tokenExpiryTime = expireTime

	return cachedToken, nil
}

// SendNotification envía una notificación usando el token de acceso cacheado.
func SendFirebaseNotification(deviceToken, title, body string) error {
	// Obtener el token de acceso (cacheado o nuevo)
	accessToken, err := GetAccessTokenWithCache()
	if err != nil {
		return err
	}

	// Obtener el TokenRefresher para usar el cliente
	tokenRefresher, err := repositories.GetFirebaseClient("gd-notificacionesandroid-firebase-adminsdk-2v5rt-090a3f0a89.json")
	if err != nil {
		return err
	}

	// Crear y enviar la solicitud de notificación
	err = repositories.SendFirebaseNotificationRequest(tokenRefresher.Client, accessToken, deviceToken, title, body, tokenRefresher.ProjectID)
	if err != nil {
		return err
	}

	return nil
}
