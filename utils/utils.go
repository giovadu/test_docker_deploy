package utils

import (
	"fmt"
	"notification_server/models"
	"strings"

	"firebase.google.com/go/v4/messaging"
)

func GenerateMessages(events []models.Events) [][]*messaging.Message {
	var messages []*messaging.Message
	for _, event := range events {
		alert := fmt.Sprintf("Alerta Vehículo %s", event.Plate)
		// Dividir los tokens por comas
		tokens := strings.Split(event.Tokens, ",")
		// Recorrer los tokens y crear un mensaje para cada uno
		for _, token := range tokens {
			// Eliminar espacios en blanco alrededor del token
			token = strings.TrimSpace(token)

			// Crear un mensaje solo si el token no está vacío
			if token != "" {
				message := &messaging.Message{
					Notification: &messaging.Notification{
						Title: alert,
						Body:  event.Event,
					},
					Token: token,
				}
				messages = append(messages, message)
			}
		}
	}
	totalMessages := SplitMessagesIntoChunks(messages, 500)
	return totalMessages
}
func SplitMessagesIntoChunks(messages []*messaging.Message, chunkSize int) [][]*messaging.Message {
	var chunks [][]*messaging.Message
	for i := 0; i < len(messages); i += chunkSize {
		end := i + chunkSize
		if end > len(messages) {
			end = len(messages)
		}
		chunks = append(chunks, messages[i:end])
	}
	return chunks
}
