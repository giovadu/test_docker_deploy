package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"notification_server/app_services"
	"notification_server/models"
	"strings"
)

func UpdateVerificationEvents(events []models.Events) error {
	var ids []int
	for i := 0; i < len(events); i++ {
		ids = append(ids, events[i].ID)
	}

	// Verificar si la lista de IDs está vacía
	if len(ids) == 0 {
		log.Println("No hay eventos para actualizar.")
		return nil // No hacer nada si no hay IDs para actualizar
	}

	db := app_services.GetConnection()

	// Construir la consulta de actualización
	query := `UPDATE tc_events SET sended = 1, notified = NOW() WHERE Id IN (`
	vals := []interface{}{}

	// Construir los placeholders y recolectar los valores
	for i, id := range ids {
		if i > 0 {
			query += ", "
		}
		query += "?"
		vals = append(vals, id)
	}

	query += ")"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparando la consulta: %v", err)
	}
	defer stmt.Close()

	// Ejecutar la consulta en una goroutine

	_, err = stmt.Exec(vals...)

	if err != nil {
		log.Printf("error ejecutando la consulta UpdateVerificationEvents: %v", err)
	}

	return nil
}

func BatchDeleteTokens(tokens []string) error {
	var tokens_parsed []string
	for i := 0; i < len(tokens); i++ {
		if !strings.Contains(tokens[i], ",") {
			tokens_parsed = append(tokens_parsed, tokens[i])
		}
	}
	if len(tokens) == 0 {
		return nil
	}
	db := app_services.GetConnection()
	// Si no hay tokens para eliminar, retornar inmediatamente

	// Construir la consulta de eliminación
	query := "DELETE FROM tc_tokens WHERE token IN ("
	placeholders := ""
	vals := []interface{}{}

	for i, tokenID := range tokens_parsed {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		vals = append(vals, tokenID)
	}

	query += placeholders + ")"

	// Ejecutar la consulta
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparando la consulta: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		return fmt.Errorf("error ejecutando la consulta: %v", err)
	}
	log.Println("Se eliminaron", len(vals), "tokens")
	return nil
}

func BatchInsertVerificationMessages(messageStatus []models.MessageStatusResponse) error {
	db := app_services.GetConnection()

	// Definir el tamaño del lote, puedes ajustarlo según el límite de tu servidor MySQL
	const batchSize = 500 // Número de filas por lote
	var vals []interface{}
	placeholders := ""
	insertCount := 0

	for i := 0; i < len(messageStatus); i++ {

		if insertCount > 0 {
			placeholders += ", "
		}
		placeholders += "(?, ?, ?)"
		vals = append(vals, messageStatus[i].AppName, messageStatus[i].Message, messageStatus[i].Token) // Usar el mensaje asociado al evento y el token

		insertCount++

		// Cuando alcanzamos el batchSize o es el último evento, ejecutar la inserción
		if insertCount >= batchSize {
			// Ejecutar la consulta de inserción en lotes
			err := executeBatchInsert(db, placeholders, vals)
			if err != nil {
				return err
			}

			// Reiniciar para el siguiente lote
			placeholders = ""
			vals = []interface{}{}
			insertCount = 0
		}

	}

	// Insertar los valores restantes que no alcanzaron a formar un lote completo
	if insertCount > 0 {
		err := executeBatchInsert(db, placeholders, vals)
		if err != nil {
			return err
		}
	}

	return nil
}

// Función auxiliar para ejecutar la consulta de inserción
func executeBatchInsert(db *sql.DB, placeholders string, vals []interface{}) error {
	query := `INSERT INTO tc_notifications_sended (appname, message, token) VALUES ` + placeholders

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparando la consulta: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		return fmt.Errorf("error ejecutando la consulta: %v", err)
	}

	return nil
}
