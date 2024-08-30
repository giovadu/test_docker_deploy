package main

import (
	"test_dcoker_deploy/app"
	"test_dcoker_deploy/http/routes"
)

func main() {
	server := app.NewServer()
	router := routes.Router()
	server.Initialize(router)
	server.Run()
}
