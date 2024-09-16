package app_services

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var clientv1 *messaging.Client
var clientv2 *messaging.Client

func InitFirebase() {
	clientv1 = initFirebasev1()
	clientv2 = initFirebasev2()
}

func GetFirebaseClientv1() *messaging.Client {
	return clientv1
}
func GetFirebaseClientv2() *messaging.Client {
	return clientv2
}
func initFirebasev1() *messaging.Client {
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
func initFirebasev2() *messaging.Client {
	fmt.Println("Initializing Firebasev2...")
	opt := option.WithCredentialsFile("notificaciones-push-1af7d-firebase-adminsdk-u0rd4-62ea0a3f0a.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		panic(err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		fmt.Printf("error initializing messaging client v2: %v", err)
		panic(err)
	}
	fmt.Println("Successfully initialized Firebase v2")
	return client
}
