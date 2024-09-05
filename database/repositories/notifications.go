package repositories

import (
	"fmt"
	"log"
	"notifcations_server/app"
)

func CreateVerificationEvent(appname, message string) error {
	db := app.GetConnection()

	_, err := db.Exec(`
	INSERT INTO tc_notifications_sended (appname, message) 
	VALUES (?, ?)`, appname, message)

	if err != nil {
		log.Printf("Error al insertar en la base de datos: %v", err)
		return err
	}
	return nil
}
func BatchDeleteTokens(tokenIDs []string) error {
	db := app.GetConnection()
	// Si no hay tokens para eliminar, retornar inmediatamente
	if len(tokenIDs) == 0 {
		return nil
	}

	// Construir la consulta de eliminación
	query := "DELETE FROM tc_tokens WHERE token IN ("
	placeholders := ""
	vals := []interface{}{}

	for i, tokenID := range tokenIDs {
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

	return nil
}
func BatchInsertVerificationMessages(appname string, messages []string) error {
	db := app.GetConnection()
	// Construir la consulta de inserción
	query := `INSERT INTO tc_notifications_sended (appname, message) VALUES `
	vals := []interface{}{}

	// Construir los placeholders y recolectar los valores
	placeholders := ""
	for i, msg := range messages {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "(?, ?)"
		vals = append(vals, appname, msg)
	}

	query += placeholders

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

	return nil
}

// func UpdateVerificationEvent(event_id int) error {
// 	db := app.GetConnection()

// 	_, err := db.Exec(`
// 	UPDATE tc_events
// 	SET sended = 1, notified = NOW()
// 	WHERE Id = ?`, event_id)

// 	if err != nil {
// 		log.Printf("Error al tc_events en la base de datos: %v", err)
// 		return err
// 	}
// 	return nil
// }
