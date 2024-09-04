package main

import (
	"notifcations_server/app"
	"notifcations_server/http/routes"
)

func main() {
	app.LoadEnv()
	app.InitMySQL()
	server := app.NewServer()
	router := routes.Router()
	server.Initialize(router)
	server.Run()
}
