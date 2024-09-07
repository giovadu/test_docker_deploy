package main

import (
	"log"
	"notification_server/app_services"
	"notification_server/models"
	"notification_server/repositories"
	"notification_server/utils"
	"sync"
	"time"
)

var StartID int = 0
var failedTokens []string
var Messages []string
var EventsReported []models.Events

func main() {
	app_services.LoadEnv()
	app_services.InitMySQL()
	app_services.InitFirebase()

	const numWorkers int = 10
	const batchSize = 500

	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex para proteger el acceso a las variables compartidas

	startID, failedTokensAux, ResponseAux, EventsReportedAux, err := processEventWithOutLimit("inicial", batchSize)
	if err == nil {
		StartID = startID
		failedTokens = append(failedTokens, failedTokensAux...)
		Messages = append(Messages, ResponseAux...)
		EventsReported = append(EventsReported, EventsReportedAux...)
	}
	for {
		startTime := time.Now()

		log.Println("Iniciando proceso de envío de notificaciones en id ", StartID)
		for i := 0; i < numWorkers; i++ {
			minID := StartID + i*batchSize
			maxID := minID + batchSize - 1
			wg.Add(1)
			go func(workerName int, minID, maxID int) {
				defer wg.Done()
				startID, failedTokensAux, ResponseAux, EventsReportedAux, err := processEventsWithLimit(workerName, minID, maxID)
				if err == nil {
					mu.Lock()
					failedTokens = append(failedTokens, failedTokensAux...)
					Messages = append(Messages, ResponseAux...)
					EventsReported = append(EventsReported, EventsReportedAux...)
					if workerName == numWorkers {
						StartID = startID
						log.Println("Worker", workerName, "actualizó el startID a", startID)
					}
					mu.Unlock()
				}

			}(i+1, minID, maxID)
		}
		wg.Wait()

		timeSince := time.Since(startTime)

		log.Printf("Inital ID: %d Este proceso demoró: %v envios totales %v, fallos totales %v ", StartID, timeSince, len(Messages), len(failedTokens))
		//se guardar registro de las notificaciones exitosas
		// repositories.BatchInsertVerificationMessages("gpsec", Messages, EventsReported)
		//se eliminan los tokens que fallaron
		repositories.BatchDeleteTokens(failedTokens)
		//se cambian de estado los eventos que se enviaron
		// repositories.UpdateVerificationEvents(EventsReported)
		failedTokens = []string{}
		Messages = []string{}
		EventsReported = []models.Events{}
	}
}

func processEventWithOutLimit(workerID string, batchSize int) (int, []string, []string, []models.Events, error) {
	eventRepom, lastId, err := repositories.GetEventsWithOutstartID(batchSize)
	if err != nil || len(eventRepom) == 0 {
		log.Printf("[Worker %s] Error obteniendo eventos: %v", workerID, err)
		return 0, []string{}, []string{}, []models.Events{}, err
	}

	messags := utils.GenetateMessages(eventRepom)

	failedTokens, Response, err := repositories.SengMessage(messags, eventRepom)
	if err != nil {
		log.Printf("[Worker %s] Error enviandos eventos: %v", workerID, err)
		return 0, failedTokens, Response, []models.Events{}, err
	}

	log.Printf("[Worker %s] Terminó de procesar eventos", workerID)
	// wg.Done()
	return lastId, failedTokens, Response, eventRepom, nil
}

func processEventsWithLimit(workerID int, minID int, maxID int) (int, []string, []string, []models.Events, error) {
	eventRepom, lastId, err := repositories.GetEventsWithLimit(minID, maxID)
	if err != nil {
		log.Printf("[Worker %d] Error obteniendo eventos: %v", workerID, err)
		// Si hay un error, retornar el mismo valor de startID
		return minID, []string{}, []string{}, []models.Events{}, err
	}

	if len(eventRepom) == 0 {
		// Si no hay más eventos, terminar el worker
		return maxID, []string{}, []string{}, []models.Events{}, nil
	}

	messags := utils.GenetateMessages(eventRepom)

	failedTokens, Response, err := repositories.SengMessage(messags, eventRepom)
	if err != nil {
		log.Printf("[Worker %d] Error enviandos eventos: %v", workerID, err)
		return 0, failedTokens, Response, []models.Events{}, err
	}

	log.Printf("[Worker %d] Terminó de procesar eventos", workerID)

	return lastId, failedTokens, Response, eventRepom, nil
}
