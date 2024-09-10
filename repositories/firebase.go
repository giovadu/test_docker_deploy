package repositories

import (
	"context"
	"notification_server/app_services"

	"firebase.google.com/go/v4/messaging"
)

func SendMessage(messages []*messaging.Message) (*messaging.BatchResponse, error) {
	firebaseClient := app_services.GetFirebaseClient()
	response, err := firebaseClient.SendEach(context.Background(), messages)
	if err != nil {
		return nil, err
	}
	return response, nil
}
