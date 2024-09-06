package services

import (
	"log"
	"notifcations_server/database/repositories"
	"sync"
	"time"
)

// Global variables to store token and expiration
var (
	cachedTokenV1     string
	tokenExpiryTimeV1 time.Time
	tokenLockV1       sync.Mutex
)
var (
	cachedTokenV2     string
	tokenExpiryTimeV2 time.Time
	tokenLockV2       sync.Mutex
)

// GetAccessTokenWithCache obtiene el token de acceso, usando uno cacheado si es válido.
func GetAccessTokenWithCacheV1() (string, error) {
	tokenLockV1.Lock()
	defer tokenLockV1.Unlock()

	// Verificar si el token ha expirado o no existe
	if time.Now().Before(tokenExpiryTimeV1) && cachedTokenV1 != "" {
		return cachedTokenV1, nil
	}
	log.Println("Creando un nuevo token de acceso para v1")
	// Obtener el TokenRefresher que maneja el refresco del token
	tokenRefresher, err := repositories.GetFirebaseClientV1("gd-notificacionesandroid-firebase-adminsdk-2v5rt-090a3f0a89.json")
	if err != nil {
		return "", err
	}

	// Obtener el nuevo token de acceso
	accessToken, err := tokenRefresher.GetAccessTokenV1()
	if err != nil {
		return "", err
	}

	// Actualizar el token y su tiempo de expiración
	expireTime := time.Now().Add(time.Minute * 50) // Ajustar según la duración real del token
	cachedTokenV1 = accessToken
	tokenExpiryTimeV1 = expireTime

	return cachedTokenV1, nil
}
func GetAccessTokenWithCacheV2() (string, error) {
	tokenLockV2.Lock()
	defer tokenLockV2.Unlock()

	// Verificar si el token ha expirado o no existe
	if time.Now().Before(tokenExpiryTimeV2) && cachedTokenV2 != "" {
		return cachedTokenV2, nil
	}
	log.Println("Creando un nuevo token de acceso para v2")
	// Obtener el TokenRefresher que maneja el refresco del token
	tokenRefresher, err := repositories.GetFirebaseClientV2("gd-notificacionesandroid-firebase-adminsdk-2v5rt-090a3f0a89.json")
	if err != nil {
		return "", err
	}

	// Obtener el nuevo token de acceso
	accessToken, err := tokenRefresher.GetAccessTokenV2()
	if err != nil {
		return "", err
	}

	// Actualizar el token y su tiempo de expiración
	expireTime := time.Now().Add(time.Minute * 50) // Ajustar según la duración real del token
	cachedTokenV2 = accessToken
	tokenExpiryTimeV2 = expireTime

	return cachedTokenV2, nil
}

// SendNotification envía una notificación usando el token de acceso cacheado.
func SendFirebaseNotificationV1(deviceToken, title, body string) error {
	// Obtener el token de acceso (cacheado o nuevo)
	accessToken, err := GetAccessTokenWithCacheV1()
	if err != nil {
		return err
	}

	// Obtener el TokenRefresher para usar el cliente
	tokenRefresher, err := repositories.GetFirebaseClientV1("gd-notificacionesandroid-firebase-adminsdk-2v5rt-090a3f0a89.json")
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
func SendFirebaseNotificationV2(deviceToken, title, body string) error {
	// Obtener el token de acceso (cacheado o nuevo)
	accessToken, err := GetAccessTokenWithCacheV2()
	if err != nil {
		return err
	}

	// Obtener el TokenRefresher para usar el cliente
	tokenRefresher, err := repositories.GetFirebaseClientV2("notificaciones-push-1af7d-firebase-adminsdk-u0rd4-62ea0a3f0a.json")
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
