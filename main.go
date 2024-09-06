package main

import (
	"fmt"
	"log"
	"notifcations_server/app"
	"notifcations_server/database/models"
	"notifcations_server/database/repositories"
	"notifcations_server/http/routes"
	"notifcations_server/services"
	"sync"
	"time"
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
	services.GetAccessTokenWithCache()
	server.Initialize(router)
	server.Run()

}
func proces() {

	startTime := time.Now() // Tiempo de inicio del ciclo
	const numWorkers = 3

	var wg sync.WaitGroup
	var mu sync.RWMutex

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			totalAux, succAux, failAux := processEvents(&mu)
			totalTime := time.Since(startTime)
			log.Printf("Ciclo completado! eventos procesados en %v, %d eventos procesados, %d exitosos, %d fallidos", totalTime, totalAux, succAux, failAux)
			wg.Done()
		}(i)
	}
	wg.Wait()

}
func processEvents(mu *sync.RWMutex) (int, int, int) {
	mu.Lock()
	eventRepom, err := repositories.GetEventsWithLimit(minID, 500)
	if len(eventRepom) > 0 {
		minID = eventRepom[len(eventRepom)-1].ID // Asume que los eventos están ordenados por ID
		log.Printf("[Worker cambió el id a %d ", minID)
	}
	mu.Unlock()
	if err != nil {
		log.Printf("[Worker error obteniendo eventos: %v", err)
		return 0, 0, 0
	}

	var wg sync.WaitGroup
	for i := 0; i < len(eventRepom); i++ {
		wg.Add(1)
		go func(event models.Events) {
			defer wg.Done()
			GetUsersToSendNotifications(event, mu)
		}(eventRepom[i])
	}

	wg.Wait()

	// Insertar todos los mensajes de verificación en la base de datos
	if len(verificationMessages) > 0 {
		go repositories.BatchInsertVerificationMessages("gpsec", verificationMessages)
	}
	if len(unregisteredTokens) > 0 {
		go repositories.BatchDeleteTokens(unregisteredTokens)
	}

	totalEvents := successCount + failCount
	return totalEvents, successCount, failCount
}

func sendNotification(user models.UserInfo, eventRepom models.Events, mu *sync.RWMutex) {
	alert := fmt.Sprintf("Alerta Vehículo %s", eventRepom.Plate)

	err := services.SendFirebaseNotification(user.Token, alert, eventRepom.Event)
	var eventMessage string

	mu.RLock()
	defer mu.RUnlock()

	if err != nil {
		unregisteredTokens = append(unregisteredTokens, user.Token)
		eventMessage = fmt.Sprintf("%s %s %d Enviando push a %s %s Dispositivo: %s error: %v", eventRepom.ServerTime, user.AppName, eventRepom.ID, user.Name, user.So, eventRepom.Plate, err)
		failCount++
	} else {
		eventMessage = fmt.Sprintf("%s %s %d Enviando push a %s %s Dispositivo: %s exitoso", eventRepom.ServerTime, user.AppName, eventRepom.ID, user.Name, user.So, eventRepom.Plate)
		successCount++
	}

	// Acumular los mensajes de verificación
	verificationMessages = append(verificationMessages, eventMessage)
}

func GetUsersToSendNotifications(eventRepom models.Events, mu *sync.RWMutex) {
	users, err := repositories.GetUsersToSendNotifications(eventRepom.Plate)
	if err != nil {
		log.Printf("error obteniendo usuarios de vehiculo %s: %v", eventRepom.Plate, err)
		return
	}

	var wg sync.WaitGroup
	for j := 0; j < len(users); j++ {
		wg.Add(1)
		go func(user models.UserInfo) {
			defer wg.Done()
			sendNotification(user, eventRepom, mu)
		}(users[j])
	}

	wg.Wait()
}
