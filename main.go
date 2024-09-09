package main

import (
	"notifcations_server/app"
	"notifcations_server/http/routes"
	"notifcations_server/services"
)

func main() {
	server := app.NewServer()
	router := routes.Router()
	services.InitFirebaseApi()
	app.InitFirebaseApp()
	server.Initialize(router)
	server.Run()
}
