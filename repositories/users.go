package repositories

import (
	"fmt"
	"notification_server/app_services"
	"notification_server/models"
	"strings"
)

func UpdateVerificationEvents(events []models.Events) error {
	var ids []int
	for i := 0; i < len(events); i++ {
		ids = append(ids, events[i].ID)
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

	return nil
}
func BatchInsertVerificationMessages(appname string, messages []string, events []models.Events) error {
	db := app_services.GetConnection()
	// Construir la consulta de inserción
	query := `INSERT INTO tc_notifications_sended (appname, message, token) VALUES `
	vals := []interface{}{}

	// Construir los placeholders y recolectar los valores
	placeholders := ""
	for i, msg := range messages {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "(?, ?, ?)"
		vals = append(vals, appname, msg, events[i].Tokens)
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
