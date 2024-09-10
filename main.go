package main

import (
	"log"
	"notification_server/app_services"
	"notification_server/repositories"
	"notification_server/utils"
	"time"

	"firebase.google.com/go/v4/messaging"
)

//TODO: VALIDAR QUE SI LA ALERTA NO TIENE EQUIVALENCIA EN TC_NOTIFICATION_EQUIVALES NO SE ENVIAR PERO SE GUARDA EN EL LOG
//TODO: VALIAR LAS CONFIGURACIONES DE ALERTAS DE LOS USUARIOS
//TODO: VALIDAR SI EL DISPOSITIVO ESTA INACTIVO
//TODO: AGREGAR LA TRADICCION DE LOS MENSAJES PARA ALARM

//FORMATO PARA EL LOG:  FECHA/HORA/APPNAME/ID DEL EVENTO/MENSAJE DE SI SE VALIDO/NOMBRE DEL USUARIO/PLACA
//ENVIO EXITOSO
//2024-09-08 00:00:00  GPSEC  5955055 Envía PUSH a Darkis Tainni Sotaban Rodriguez (Apple iOS 16.1.1 IPHONE) Dispositivo SFY91G: Vehículo ha sido apagado  velocidad: 8 Km/h RESPONSE: 200

//EJEMPLO ENVIO SIN EQUIVALENCIA
//2024-09-08 00:00:00  GPSEC  5955028  NO se envía PUSH a Julián Andrés Jiménez Cabanzo pues no se encuentra la configuración de notificación para deviceStopped Dispositivo ZMA86G:

// EJEMPLO ENVIO CUANDO USUARIO TIENE DESACTIVADAS LAS ALERTAS
// 2024-09-08 00:00:05  GPSEC  5955153  NO se envía PUSH a Nicolás Santiago Castro Zambrano por estar desactivado para este usuario Dispositivo DGZ16G: Vehículo ha excedido el límite de velocidad  velocidad: 74 Km/h
func main() {
	init_counter := 0
	app_services.LoadEnv()
	app_services.InitMySQL()
	app_services.InitFirebase()
	const batchSize = 500
	init_counter, _, err := handleMessages(batchSize, init_counter)
	if err != nil || init_counter == 0 {
		log.Panic("error iniciando el programa: %v", err)
		return
	}

	//batch 500 =>  eventos 48/s, notificacions  86/s, tiempo del test 73s
	//batch 1000 => eventos 70/s, notificacions 110/s, tiempo del test 71s
	//batch 1500 => eventos 82/s, notificacions 131/s, tiempo del test 73s

	for {
		init_counterAux, lenght, err := handleMessages(batchSize, init_counter)
		if err != nil {
			log.Println("error en el proceso de envío: %v", err)
			return
		}
		if init_counterAux == 0 {
			time.Sleep(1 * time.Second)
		} else if lenght < 100 {
			time.Sleep(1 * time.Second)
		} else {
			init_counter = init_counterAux
		}
	}
}
func handleMessages(batchSize int, init_counter int) (int, int, error) {
	startTime := time.Now()
	log.Println("Iniciando proceso de envio de mensajes")
	eventRepom, err := repositories.GetEventsWithOutstartID(init_counter, batchSize)
	if err != nil || len(eventRepom) == 0 {
		log.Printf("[Worker inicial] Error obteniendo eventos: %v", err)
		return 0, 0, err
	}

	go func() {
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
	}()
	return eventRepom[len(eventRepom)-1].ID, len(eventRepom), nil

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
