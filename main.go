package main

import (
	"fmt"
	"log"
	"notifcations_server/app"
	"notifcations_server/database/models"
	"notifcations_server/database/repositories"
	"notifcations_server/services"
	"sync"
	"time"
)

var startID int = 0

func main() {
	app.LoadEnv() // Cargar los datos del .env
	app.InitMySQL()
	processEvents(startID)
}

func processEvents(id int) {
	startTime := time.Now() // Tiempo de inicio del ciclo
	eventRepom, err := repositories.GetEventsWithLimit(id, 700)
	if err != nil {
		log.Printf("[Worker %d] error obteniendo eventos: %v", id, err)
		return
	}

	var wgEvents sync.WaitGroup
	var successCount, failCount int

	var verificationMessages []string
	var unregisteredTokens []string

	for i := 0; i < len(eventRepom); i++ {
		wgEvents.Add(1) // Añadir un contador al WaitGroup para cada evento
		go func(event models.Events) {
			defer wgEvents.Done() // Decrementa el contador cuando la goroutine finaliza

			users, err := repositories.GetUsersToSendNotifications(event.Plate)
			if err != nil {
				log.Printf("[Worker %d] error obteniendo usuarios: %v", id, err)
				return
			}

			var wgUsers sync.WaitGroup
			for j := 0; j < len(users); j++ {
				wgUsers.Add(1)

				defer wgUsers.Done()

				alert := fmt.Sprintf("Alerta Vehículo %s", event.Plate)

				err := services.SendFirebaseNotification(users[j].Token, alert, event.Event)
				var eventMessage string
				if err != nil {
					unregisteredTokens = append(unregisteredTokens, users[j].Token)
					eventMessage = fmt.Sprintf("%s %s %d Enviando push a %s %s Dispositivo: %s error: %v", event.ServerTime, users[j].AppName, event.ID, users[j].Name, users[j].So, event.Plate, err)
					failCount++
				} else {
					eventMessage = fmt.Sprintf("%s %s %d Enviando push a %s %s Dispositivo: %s exitoso", event.ServerTime, users[j].AppName, event.ID, users[j].Name, users[j].So, event.Plate)
					successCount++
				}

				// Acumular los mensajes de verificación
				verificationMessages = append(verificationMessages, eventMessage)

			}

			// Espera a que todas las goroutines de usuarios terminen
			wgUsers.Wait()
		}(eventRepom[i])
	}

	// Espera a que todas las goroutines de eventos terminen
	wgEvents.Wait()

	// Insertar todos los mensajes de verificación en la base de datos
	if len(verificationMessages) > 0 {
		repositories.BatchInsertVerificationMessages("gpsec", verificationMessages)
	}
	if len(unregisteredTokens) > 0 {
		err := repositories.BatchDeleteTokens(unregisteredTokens)
		if err != nil {
			log.Printf("Total eliminados: %d  ", len(unregisteredTokens))
		}
	}

	totalEvents := successCount + failCount
	// Tiempo de fin del ciclo
	totalTime := time.Since(startTime)

	// Imprimir estadísticas del ciclo
	log.Printf("[Worker %d] Ciclo completado: %d eventos procesados en %v", id, totalEvents, totalTime)
	log.Printf("[Worker %d] Ciclo completado: %d eventos procesados, %d exitosos, %d fallidos", id, totalEvents, successCount, failCount)
}

//1010 eventos procesados en 38.79007075s
