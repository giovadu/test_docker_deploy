package app_services

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var client *messaging.Client

func InitFirebase() {
	client = initFirebase()
}
func GetFirebaseClient() *messaging.Client {
	return client
}

func initFirebase() *messaging.Client {
	fmt.Println("Initializing Firebase...")
	opt := option.WithCredentialsFile("gd-notificacionesandroid-firebase-adminsdk-2v5rt-090a3f0a89.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		panic(err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		fmt.Printf("error initializing messaging client: %v", err)
		panic(err)
	}
	fmt.Println("Successfully initialized Firebase")
	return client
}
