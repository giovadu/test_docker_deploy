package main

import (
	"notifcations_server/app"
	"notifcations_server/http/routes"
	"notifcations_server/services"
)

var minID int = 0

var (
	verificationMessages []string
	unregisteredTokens   []string
	successCount         int
	failCount            int
)

func main() {

	server := app.NewServer()
	router := routes.Router()
	services.GetAccessTokenWithCacheV1()
	services.GetAccessTokenWithCacheV2()
	server.Initialize(router)
	server.Run()

}
