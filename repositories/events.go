package repositories

import (
	"context"
	"fmt"
	"notification_server/app_services"
	"notification_server/models"
)

func GetEventsWithOutstartID(startID, limit int) (events []models.Events, err error) {
	// Inicializa el slice de eventos
	events = make([]models.Events, 0)

	// Obtén la conexión a la base de datos
	db := app_services.GetConnection()

	// rows, err := db.Query("SELECT * FROM tracker.tc_events_priority_1 LIMIT ?", limit)
	rows, err := db.QueryContext(context.Background(), `
        SELECT 
            e.id AS id,
            e.type AS type,
            COALESCE(CONCAT(n.translate,
                            IF(JSON_EXTRACT(e.attributes, '$.speed') <> '',
                                CONCAT(' a ',
                                        ROUND(JSON_EXTRACT(e.attributes, '$.speed') * 1.852, 0),
                                        ' Km/h'),
                                '')),
                    'Alarma') AS event,
            d.appname AS appname,
            d.name AS name,
            u.name AS username,
            u.tokens AS tokens,
            COALESCE(g.name, '') AS geofencename
        FROM
            tracker.tc_events e
            JOIN tracker.tc_devices d ON d.id = e.deviceid
            LEFT JOIN tracker.tc_geofences g ON e.geofenceid = g.id
            JOIN tracker.tc_notification_locate n ON e.type = n.original
            JOIN tracker.tc_users u ON FIND_IN_SET(d.name, u.devices) > 0
        WHERE
            u.tokens <> ''
            AND n.translate IS NOT NULL
            AND e.sended = 0
			AND e.eventtime > (NOW() - INTERVAL 1 MINUTE)
			AND e.id > ?
        LIMIT ?`, startID, limit)

	if err != nil {
		return nil, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var event models.Events
		if err := rows.Scan(&event.ID, &event.Type, &event.Event, &event.AppName, &event.Plate, &event.UserName, &event.Tokens, &event.GeofenceName); err != nil {
			return nil, fmt.Errorf("error al preparar la consulta: %w", err)
		}
		if event.Type == "geofenceExit" || event.Type == "geofenceEnter" {
			event.Event = fmt.Sprintf("%s %s", event.Event, event.GeofenceName)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de resultados: %w", err)
	}

	return events, nil
}
