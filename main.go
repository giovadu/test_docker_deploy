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

	errorChannel := make(chan error, len(eventRepom)) // Canal para capturar errores
	var successCount, failCount int

	// Slice para acumular los mensajes de verificación
	var verificationMessages []string
	var unregisteredTokens []string

	for i := 0; i < len(eventRepom); i++ {
		wgEvents.Add(1) // Añadir un contador al WaitGroup para cada evento
		go func(event models.Events) {
			defer wgEvents.Done() // Decrementa el contador cuando la goroutine finaliza

			users, err := repositories.GetUsersToSendNotifications(event.Plate)
			if err != nil {
				errorChannel <- fmt.Errorf("error obteniendo usuarios para %s: %v", event.Plate, err)
				failCount++
				return
			}

			var wgUsers sync.WaitGroup
			for j := 0; j < len(users); j++ {
				wgUsers.Add(1)
				go func(user models.UserInfo) {
					defer wgUsers.Done()

					alert := fmt.Sprintf("Alerta Vehículo %s", event.Plate)

					err := services.SendFirebaseNotification(user.Token, alert, event.Event)
					var eventMessage string
					if err != nil {
						unregisteredTokens = append(unregisteredTokens, user.Token)
						eventMessage = fmt.Sprintf("%s %s %d Enviando push a %s %s Dispositivo: %s error: %v", event.ServerTime, user.AppName, event.ID, user.Name, user.So, event.Plate, err)
						errorChannel <- err
						failCount++
					} else {
						eventMessage = fmt.Sprintf("%s %s %d Enviando push a %s %s Dispositivo: %s exitoso", event.ServerTime, user.AppName, event.ID, user.Name, user.So, event.Plate)
						successCount++
					}

					// Acumular los mensajes de verificación
					verificationMessages = append(verificationMessages, eventMessage)

				}(users[j])
			}

			// Espera a que todas las goroutines de usuarios terminen
			wgUsers.Wait()
		}(eventRepom[i])
	}

	// Espera a que todas las goroutines de eventos terminen
	wgEvents.Wait()
	close(errorChannel)

	// Maneja los errores capturados
	log.Println(len(errorChannel))

	// Insertar todos los mensajes de verificación en la base de datos
	if len(verificationMessages) > 0 {
		go repositories.BatchInsertVerificationMessages("gpsec", verificationMessages)
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
