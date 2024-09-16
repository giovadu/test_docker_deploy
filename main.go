package main

import (
	"log"
	"notification_server/app_services"
	"notification_server/models"
	"notification_server/repositories"
	"notification_server/utils"
	"time"

	"firebase.google.com/go/v4/messaging"
	env "github.com/joho/godotenv"
)

func LoadEnv() {
	err := env.Load(".env")
	if err != nil {
		panic(err)
	}
}
func main() {
	LoadEnv()
	init_counter := 0
	app_services.InitMySQL()
	app_services.InitFirebase()

	// Ahora ejecutar cada 8 horas en una goroutine
	go func() {
		for {
			log.Println("Ejecutando actualización periódica de eventos y equivalencias de appname...")
			repositories.GetEventsTranslated()
			repositories.GetAppnameEquivalent()
			log.Println("Actualización periódica completa. Esperando 8 horas para la próxima ejecución.")
			time.Sleep(24 * time.Hour) // Espera de 8 horas
		}
	}()

	const batchSize = 1000
	var err error
	init_counter, _, err = handleMessages(batchSize, init_counter)
	if err != nil || init_counter == 0 {
		log.Printf("error iniciando el programa: %v", err)
		return
	}

	for {
		init_counterAux, count, err := handleMessages(batchSize, init_counter)
		if err != nil {
			time.Sleep(3 * time.Second)
			log.Printf("error en el proceso de envío: %v", err)
			continue
		}
		if init_counterAux != 0 && init_counterAux > init_counter {
			log.Println("Se enviaron", count, "cambió el id de:", init_counter, " a ", init_counterAux)
			init_counter = init_counterAux

		} else {
			time.Sleep(2 * time.Second)
		}
	}
}
func handleMessages(batchSize int, init_counter int) (int, int, error) {
	// startTime := time.Now()
	// log.Println("Iniciando proceso de envio de mensajes")
	eventRepom, err := repositories.GetEventsWithOutstartID(init_counter, batchSize)
	if err != nil {
		log.Printf("[Worker inicial] Error obteniendo eventos: %v", err)
		return 0, 0, err
	}
	if len(eventRepom) == 0 {
		log.Println("No hay eventos para enviar")
		return 0, 0, nil
	}
	go func() {
		//se generan los mensaje a enviar
		messages, failedMessages := utils.GenerateMessages(eventRepom)
		//se crea una variable para almacenar los valores respondidos
		var total_message_sended_in_batchs []*messaging.SendResponse

		var messages_to_compare []models.MessageStatus
		//ciclo para enviar mensajes por batches
		for i := 0; i < len(messages); i++ {
			messages_to_compare = append(messages_to_compare, messages[i]...)
			var messages_to_send_v1 []*messaging.Message
			var messages_to_send_v2 []*messaging.Message
			for j := 0; j < len(messages[i]); j++ {
				if messages[i][j].Event.Equivalent == "v1" {
					messages_to_send_v1 = append(messages_to_send_v1, messages[i][j].Message)
				} else {
					messages_to_send_v2 = append(messages_to_send_v2, messages[i][j].Message)
				}
			}
			if len(messages_to_send_v1) != 0 {
				BatchResponse, err := repositories.SendMessageV1(messages_to_send_v1)
				if err != nil {
					log.Printf("Error enviando mensajes: %v", err)
					return
				}
				if len(BatchResponse.Responses) != 0 {
					total_message_sended_in_batchs = append(total_message_sended_in_batchs, BatchResponse.Responses...)
				}
			}
			if len(messages_to_send_v2) != 0 {
				BatchResponseV2, err := repositories.SendMessageV2(messages_to_send_v2)
				if err != nil {
					log.Printf("Error enviando mensajes: %v", err)
					return
				}

				if len(BatchResponseV2.Responses) != 0 {
					total_message_sended_in_batchs = append(total_message_sended_in_batchs, BatchResponseV2.Responses...)
				}
			}
		}
		//vamos a analizar la respuesta
		var final_messages []models.MessageStatusResponse
		for i := 0; i < len(total_message_sended_in_batchs); i++ {
			if total_message_sended_in_batchs[i].Success {
				final_messages = append(final_messages, models.FormatStatusMessage(messages_to_compare[i].Event, true, "STATUS 200"))

			} else {
				final_messages = append(final_messages, models.FormatStatusMessage(messages_to_compare[i].Event, false, "STATUS 400 TOKEN NO REGISTRADO"))
			}
		}
		final_messages = append(final_messages, failedMessages...)
		go repositories.BatchInsertVerificationMessages(final_messages)
	}()
	lastEventID := eventRepom[0].ID

	for _, event := range eventRepom {
		if event.ID > lastEventID {
			lastEventID = event.ID
		}
	}

	return lastEventID, len(eventRepom), nil
}
