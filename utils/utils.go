package utils

import (
	"fmt"
	"strings"

	"firebase.google.com/go/v4/messaging"
)

func GenerateMessages(alert, message, tokens string) ([]*messaging.Message, error) {
	var messages []*messaging.Message
	// Dividir los tokens en un slice
	tokensParsed := strings.Split(tokens, ",")
	if len(tokensParsed) == 0 || len(tokensParsed) >= 501 {
		return messages, fmt.Errorf("ha excecido el límite de tokens permitidos")
	}

	for _, token := range tokensParsed {
		// Eliminar espacios en blanco alrededor del token
		token = strings.TrimSpace(token)

		// Crear un mensaje solo si el token no está vacío
		if token != "" {
			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: alert,
					Body:  message,
				},
				Token: token,
				Android: &messaging.AndroidConfig{
					Priority: "high",
					Notification: &messaging.AndroidNotification{
						Sound: "default",
					},
				},
				APNS: &messaging.APNSConfig{
					Payload: &messaging.APNSPayload{
						Aps: &messaging.Aps{
							Sound: "default",
						},
					},
					Headers: map[string]string{
						"apns-priority": "10",
					},
				},
			}
			messages = append(messages, message)
		}
	}
	return messages, nil
}
