package app

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var clientV1 *messaging.Client
var clientV2 *messaging.Client

func InitFirebaseApp() {
	clientV1 = initFirebaseV1()
	clientV2 = initFirebaseV2()
}
func GetFirebaseClientV1() *messaging.Client {
	return clientV1
}
func GetFirebaseClientV2() *messaging.Client {
	return clientV2
}

func initFirebaseV1() *messaging.Client {
	fmt.Println("Initializing Firebase V1...")
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
	fmt.Println("Successfully initialized Firebase V1")
	return client
}
func initFirebaseV2() *messaging.Client {
	fmt.Println("Initializing FirebaseV2...")
	opt := option.WithCredentialsFile("notificaciones-push-1af7d-firebase-adminsdk-u0rd4-62ea0a3f0a.json")
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
	fmt.Println("Successfully initialized Firebase V2")
	return client
}
