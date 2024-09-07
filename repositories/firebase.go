package repositories

import (
	"context"
	"fmt"
	"notification_server/app_services"
	"notification_server/models"

	"firebase.google.com/go/v4/messaging"
)

func SengMessage(messages []*messaging.Message, eventRepom []models.Events) ([]string, []string, error) {
	//obtengo el cliente
	var failedTokens []string
	var statusMessages []string
	firebaseClient := app_services.GetFirebaseClient()
	//envio los 500 mensajes
	response, err := firebaseClient.SendEachDryRun(context.Background(), messages)
	//valido si envió correctamente
	if err != nil {
		fmt.Printf("error sending message: %v", err)
		return []string{}, []string{}, err
	}
	//agrego los indices de los tokens que fallaron
	for i := 0; i < len(response.Responses); i++ {
		if response.Responses[i].Error != nil {
			failedTokens = append(failedTokens, messages[i].Token)
			statusMessages = append(statusMessages, fmt.Sprintf("Error enviando mensaje a %s: %v", messages[i].Token, response.Responses[i].Error))
		} else {
			statusMessages = append(statusMessages, fmt.Sprintf("Mensaje enviado a %s", messages[i].Token))
		}
	}
	return failedTokens, statusMessages, nil
}
