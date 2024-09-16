package repositories

import (
	"context"
	"notification_server/app_services"

	"firebase.google.com/go/v4/messaging"
)

func SendMessageV1(messages []*messaging.Message) (*messaging.BatchResponse, error) {
	firebaseClient := app_services.GetFirebaseClientv1()
	response, err := firebaseClient.SendEach(context.Background(), messages)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func SendMessageV2(messages []*messaging.Message) (*messaging.BatchResponse, error) {
	firebaseClient := app_services.GetFirebaseClientv2()
	response, err := firebaseClient.SendEach(context.Background(), messages)
	if err != nil {
		return nil, err
	}
	return response, nil
}
