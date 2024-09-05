package repositories

import (
	"fmt"
	"notifcations_server/app"
	"notifcations_server/database/models"
)

func GetEventsWithLimit(startID int, limit int) (events []models.Events, err error) {
	// Inicializa el slice de eventos
	events = make([]models.Events, 0)

	// Obtén la conexión a la base de datos
	db := app.GetConnection()

	// Consulta SQL para obtener todos los parámetros
	query := fmt.Sprintf(`
	SELECT
	    e.id AS id,
	    p.servertime AS servertime,
	    e.type AS type,
	    COALESCE(et.translate, 'Alarma') AS event,
	    d.name as device_name,
		g.name as geofencename
	FROM
	    tracker.tc_events e
	JOIN
	    tracker.tc_devices d ON e.deviceid = d.id
	JOIN
	    tracker.tc_positions p ON p.id = e.positionid
	LEFT JOIN
	    tracker.tc_notification_locate et ON e.type = et.original
	LEFT JOIN
	    tracker.tc_notification_locate nl1 ON JSON_UNQUOTE(JSON_EXTRACT(e.attributes, '$.alarm')) = nl1.original
	LEFT JOIN
	    tracker.tc_geofences g ON e.geofenceid = g.id
	WHERE
	    e.sended = 0
	    AND p.servertime > (NOW() - INTERVAL 3 MINUTE)
	    AND e.positionid > 0
	    AND e.type <> 'commandResult'
	    AND d.active = 1
	    AND (et.translate IS NOT NULL OR nl1.translate IS NOT NULL)
		AND e.id > %d
	ORDER BY e.id
	    LIMIT %d;
    `, startID, limit)

	// Prepara la consulta
	stmt, err := db.Prepare(query)
	if err != nil {
		return events, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return events, fmt.Errorf("error ejecutando la consulta: %w", err)
	}
	defer rows.Close()
	// Itera sobre los resultados y los escanea en el slice de eventos
	for rows.Next() {
		var event models.Events
		err := rows.Scan(
			&event.ID,
			&event.ServerTime,
			&event.Type,
			&event.Event,
			&event.Plate,
			&event.GeofenceName,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear el resultado: %w", err)
		}

		if event.Type == "geofenceExit" || event.Type == "geofenceEnter" {
			event.Event = fmt.Sprintf("%s %s", event.Event, event.GeofenceName)
		}
		events = append(events, event)
	}

	// Verifica si hubo algún error durante la iteración de resultados
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de resultados: %w", err)
	}

	return events, nil
}
func GetUsersToSendNotifications(plate string) (users []models.UserInfo, err error) {
	// Inicializa el slice de eventos
	users = make([]models.UserInfo, 0)

	// Obtén la conexión a la base de datos
	db := app.GetConnection()

	// Consulta SQL para obtener todos los parámetros
	query := fmt.Sprintf(`
	SELECT
	    d.name  AS plate,
	    u.name AS username,
	    t.token AS token,
	    t.so AS so,
		d.appname as appname
	FROM
	    tracker.tc_devices d
	INNER JOIN
	    tracker.tc_users u ON FIND_IN_SET(d.name, u.devices) > 0
	INNER JOIN
	    tracker.tc_tokens t ON t.phone = u.phone
	WHERE 
		d.name = '%s';
    `, plate)

	// Prepara la consulta
	stmt, err := db.Prepare(query)
	if err != nil {
		return users, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return users, fmt.Errorf("error ejecutando la consulta: %w", err)
	}
	defer rows.Close()
	// Itera sobre los resultados y los escanea en el slice de eventos
	for rows.Next() {
		var user models.UserInfo
		err := rows.Scan(
			&user.Plate,
			&user.Name,
			&user.Token,
			&user.So,
			&user.AppName,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear el resultado: %w", err)
		}
		users = append(users, user)
	}

	// Verifica si hubo algún error durante la iteración de resultados
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de resultados: %w", err)
	}

	return users, nil
}
