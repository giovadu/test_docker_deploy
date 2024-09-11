package utils

import (
	"fmt"
	"notification_server/models"
	"strings"

	"firebase.google.com/go/v4/messaging"
)

func GenerateMessages(events []models.Events) ([][]models.MessageStatus, []models.MessageStatusResponse) {
	var messages []models.MessageStatus
	var failedMessages []models.MessageStatusResponse

	for _, event := range events {
		eventMap := models.StructToMap(event)
		if eventMap[event.Type] == 0 {
			failedMessages = append(failedMessages, models.FormatStatusMessage(event, false, "Porque en su configuración tiene desactivado este tipo de alerta"))
			continue
		}

		alert := fmt.Sprintf("Alerta Vehículo %s", event.Plate)
		tokens := strings.Split(event.Tokens, ",")
		for _, token := range tokens {
			token = strings.TrimSpace(token)
			if token != "" {
				message := models.MessageStatus{
					Message: &messaging.Message{
						Notification: &messaging.Notification{
							Title: alert,
							Body:  event.Event,
						},
						Token: token,
					},
					Event: event,
				}
				messages = append(messages, message)
			}
		}
	}

	totalMessages := splitMessagesIntoChunks(messages, 500)
	return totalMessages, failedMessages
}

func splitMessagesIntoChunks(messages []models.MessageStatus, chunkSize int) [][]models.MessageStatus {
	var chunks [][]models.MessageStatus
	for chunkSize < len(messages) {
		messages, chunks = messages[chunkSize:], append(chunks, messages[0:chunkSize:chunkSize])
	}
	return append(chunks, messages)
}
