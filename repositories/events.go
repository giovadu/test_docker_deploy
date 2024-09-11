package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"notification_server/app_services"
	"notification_server/models"
)

var Tranlates = make(map[string]string)

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
			e.attributes as attributes,
            d.appname AS appname,
            d.name AS name,
            u.name AS username,
            u.tokens AS tokens,
            COALESCE(g.name, '') AS geofencename,
			u.notification_door AS door,
			u.notification_powerOn AS ignitionOn,
			u.notification_powerOff AS ignitionOff,
			u.notification_deviceOverspeed AS deviceOverspeed,
			u.notification_geofenceEnter AS geofenceEnter,
			u.notification_geofenceExit AS geofenceExit,
			u.notification_shock AS shock,
			u.notification_powerCut AS powerCut,
			u.notification_lowbattery AS lowBattery,
			u.notification_sos AS sos
        FROM
            tracker.tc_events e
            JOIN tracker.tc_devices d ON d.id = e.deviceid
            LEFT JOIN tracker.tc_geofences g ON e.geofenceid = g.id
            JOIN tracker.tc_notification_locate n ON e.type = n.original
            JOIN tracker.tc_users u ON FIND_IN_SET(d.name, u.devices) > 0
        WHERE
            u.tokens <> ''
			AND d.active = 1
            AND n.translate IS NOT NULL
			AND e.type NOT IN ('commandResult', 'deviceMoving', 'deviceOffline', 'deviceOnline', 'deviceStopped', 'deviceUnknown')
			AND e.attributes <> '{"alarm":"fuelLeak"}'
			AND e.eventtime > (NOW() - INTERVAL 1 MINUTE)
			AND e.id > ?
		ORDER BY e.id ASC
        LIMIT ?`, startID, limit)

	if err != nil {
		return nil, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var event models.Events

		if err := rows.Scan(&event.ID,
			&event.Type,
			&event.Event,
			&event.Attributes,
			&event.AppName,
			&event.Plate,
			&event.UserName,
			&event.Tokens,
			&event.GeofenceName,
			&event.NotificationDoor,
			&event.NotificationPowerOn,
			&event.NotificationPowerOff,
			&event.NotificationDeviceOverspeed,
			&event.NotificationGeofenceEnter,
			&event.NotificationGeofenceExit,
			&event.NotificationShock,
			&event.NotificationPowerCut,
			&event.NotificationLowBattery,
			&event.NotificationSos); err != nil {
			return nil, fmt.Errorf("error al preparar la consulta: %w", err)
		}
		if event.Type == "geofenceExit" || event.Type == "geofenceEnter" {
			event.Event = fmt.Sprintf("%s %s", event.Event, event.GeofenceName)
		}
		if event.Type == "alarm" {
			var attributes map[string]string
			err := json.Unmarshal([]byte(event.Attributes), &attributes)
			if err != nil {
				fmt.Println("Error al parsear JSON:", err)
			}
			alarmType, exists := Tranlates[attributes["alarm"]]
			if exists {
				event.Type = attributes["alarm"]
				event.Event = alarmType
			}

		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de resultados: %w", err)
	}

	return events, nil
}
func GetEventsTranslated() {
	db := app_services.GetConnection()
	rows, err := db.QueryContext(context.Background(), `SELECT  original, translate  FROM tracker.tc_notification_locate `)
	if err != nil {
		panic(fmt.Errorf("error al preparar la consulta: %w", err))
	}
	defer rows.Close()
	for rows.Next() {
		var Original string
		var Translate string
		if err := rows.Scan(&Original, &Translate); err != nil {
			panic(fmt.Errorf("error al preparar la consulta: %w", err))
		}
		Tranlates[Original] = Translate
	}
	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("error durante la iteración de resultados: %w", err))
	}
}
