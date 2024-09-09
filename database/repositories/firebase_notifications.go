package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"notifcations_server/app"
	"os"
	"time"

	"firebase.google.com/go/v4/messaging"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

// TokenRefresher almacena el cliente, el TokenSource y el ProjectID.
type TokenRefresherV1 struct {
	Client      *http.Client
	TokenSource oauth2.TokenSource
	ProjectID   string
}
type TokenRefresherV2 struct {
	Client      *http.Client
	TokenSource oauth2.TokenSource
	ProjectID   string
}

// GetFirebaseClient obtiene el TokenRefresherV1 para manejar el refresco del token.
func GetFirebaseClientV1(credentialsFilePath string) (*TokenRefresherV1, error) {
	ctx := context.Background()

	// Leer el archivo JSON de credenciales
	credsData, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de credenciales: %v", err)
	}

	// Cargar las credenciales desde el archivo JSON
	creds, err := google.CredentialsFromJSON(ctx, credsData, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("error cargando credenciales: %v", err)
	}

	// Crear cliente HTTP con las credenciales
	client, _, err := transport.NewHTTPClient(ctx, option.WithTokenSource(creds.TokenSource))

	if err != nil {
		return nil, fmt.Errorf("error creando cliente HTTP: %v", err)
	}

	// Crear un TokenRefresher que maneja el refresco del token
	tokenRefresher := &TokenRefresherV1{
		Client:      client,
		TokenSource: creds.TokenSource,
		ProjectID:   creds.ProjectID,
	}

	return tokenRefresher, nil
}

// GetAccessToken obtiene el token de acceso actual. El token se refrescará si es necesario.
func (tr *TokenRefresherV1) GetAccessTokenV1() (string, error) {
	token, err := tr.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("error obteniendo el token de acceso: %v", err)
	}
	log.Printf("Nuevo token obtenido ")

	// Mostrar tiempo de expiración del token
	expireTime := time.Until(token.Expiry)
	log.Printf("El token expira a las: %v", expireTime)
	return token.AccessToken, nil
}

func GetFirebaseClientV2(credentialsFilePath string) (*TokenRefresherV2, error) {
	ctx := context.Background()

	// Leer el archivo JSON de credenciales
	credsData, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de credenciales: %v", err)
	}

	// Cargar las credenciales desde el archivo JSON
	creds, err := google.CredentialsFromJSON(ctx, credsData, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("error cargando credenciales: %v", err)
	}

	// Crear cliente HTTP con las credenciales
	client, _, err := transport.NewHTTPClient(ctx, option.WithTokenSource(creds.TokenSource))

	if err != nil {
		return nil, fmt.Errorf("error creando cliente HTTP: %v", err)
	}

	// Crear un TokenRefresher que maneja el refresco del token
	tokenRefresher := &TokenRefresherV2{
		Client:      client,
		TokenSource: creds.TokenSource,
		ProjectID:   creds.ProjectID,
	}

	return tokenRefresher, nil
}

// GetAccessToken obtiene el token de acceso actual. El token se refrescará si es necesario.
func (tr *TokenRefresherV2) GetAccessTokenV2() (string, error) {
	token, err := tr.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("error obteniendo el token de acceso: %v", err)
	}
	log.Printf("Nuevo token obtenido ")

	// Mostrar tiempo de expiración del token
	expireTime := time.Until(token.Expiry)
	log.Printf("El token expira a las: %v", expireTime)
	return token.AccessToken, nil
}

// SendNotificationRequest envía la solicitud de notificación a Firebase.
func SendFirebaseNotificationRequest(client *http.Client, accessToken, deviceToken, title, body, projectID string) error {
	msg := struct {
		Message struct {
			Token        string `json:"token"`
			Notification struct {
				Title string `json:"title"`
				Body  string `json:"body"`
			} `json:"notification"`
		} `json:"message"`
	}{
		Message: struct {
			Token        string `json:"token"`
			Notification struct {
				Title string `json:"title"`
				Body  string `json:"body"`
			} `json:"notification"`
		}{
			Token: deviceToken,
			Notification: struct {
				Title string `json:"title"`
				Body  string `json:"body"`
			}{
				Title: title,
				Body:  body,
			},
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error serializando el mensaje: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectID), bytes.NewBuffer(jsonMsg))
	if err != nil {
		return fmt.Errorf("error creando la solicitud: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando la solicitud: %v", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo la respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", string(bodyResp))
	}

	return nil
}
func SendMessagesV1(messages []*messaging.Message) (*messaging.BatchResponse, error) {
	firebaseClient := app.GetFirebaseClientV1()
	response, err := firebaseClient.SendEach(context.Background(), messages)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func SendMessagesV2(messages []*messaging.Message) (*messaging.BatchResponse, error) {
	firebaseClient := app.GetFirebaseClientV2()
	response, err := firebaseClient.SendEach(context.Background(), messages)
	if err != nil {
		return nil, err
	}
	return response, nil
}
