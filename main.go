package main

import (
	"log"
	"notification_server/app_services"
	"notification_server/repositories"
	"notification_server/utils"
	"time"

	"firebase.google.com/go/v4/messaging"
)

func main() {
	app_services.LoadEnv()
	app_services.InitMySQL()
	app_services.InitFirebase()
	const batchSize = 1500

	//batch 500 =>  eventos 48/s, notificacions  86/s, tiempo del test 73s
	//batch 1000 => eventos 70/s, notificacions 110/s, tiempo del test 71s
	//batch 1500 => eventos 82/s, notificacions 131/s, tiempo del test 73s
	for {
		startTime := time.Now()
		log.Println("Iniciando proceso de envio de mensajes")
		eventRepom, err := repositories.GetEventsWithOutstartID(batchSize, 0)
		if err != nil || len(eventRepom) == 0 {
			log.Printf("[Worker inicial] Error obteniendo eventos: %v", err)
			return
		}
		messages := utils.GenerateMessages(eventRepom)
		var statusSened []*messaging.SendResponse
		for i := 0; i < len(messages); i++ {
			BatchResponse, err := repositories.SendMessage(messages[i])
			if err != nil {
				log.Printf("Error enviando mensajes: %v", err)
				return
			}
			if len(BatchResponse.Responses) != 0 {
				statusSened = append(statusSened, BatchResponse.Responses...)
			}
		}
		totalTime := time.Since(startTime)
		log.Println("Proceso de envio de mensajes finalizado en: ", totalTime, " segundos y envió ", len(statusSened), " mensajes")
	}
}

// // Función para leer el número desde el archivo
// func readNumberFromFile(filename string) (int, error) {
// 	// Leer el archivo
// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return 0, fmt.Errorf("error al leer el archivo: %v", err)
// 	}

// 	// Convertir el contenido a un número
// 	number, err := strconv.Atoi(string(data))
// 	if err != nil {
// 		return 0, fmt.Errorf("error al convertir el contenido a un número: %v", err)
// 	}

// 	return number, nil
// }

// // Función para escribir el número en el archivo
// func writeNumberToFile(filename string, number int) error {
// 	// Convertir el número a string
// 	newData := strconv.Itoa(number)

// 	// Escribir el nuevo valor en el archivo
// 	err := ioutil.WriteFile(filename, []byte(newData), 0644)
// 	if err != nil {
// 		return fmt.Errorf("error al escribir en el archivo: %v", err)
// 	}

// 	return nil
// }

// func main() {
// 	// Ruta completa del archivo en tu escritorio (asegúrate de cambiar "tu_usuario" por tu nombre de usuario)
// 	filename := "/Users/tu_usuario/Desktop/archivo.txt"

// 	// Leer el número desde el archivo
// 	number, err := readNumberFromFile(filename)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Incrementar el número (o cualquier otra operación)
// 	number += 1

// 	// Escribir el nuevo número en el archivo
// 	err = writeNumberToFile(filename, number)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Printf("El número ha sido incrementado a: %d\n", number)
// }
