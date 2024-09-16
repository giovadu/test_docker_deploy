package models

import (
	"fmt"
	"time"

	"firebase.google.com/go/v4/messaging"
)

// Events representa la estructura de un evento en la base de datos.
type Events struct {
	ID                          int    `gorm:"column:id" json:"id"`
	Type                        string `gorm:"column:type" json:"type"`
	Event                       string `gorm:"column:event" json:"event"`
	AppName                     string `gorm:"column:appname"`
	Plate                       string `gorm:"column:name"`
	UserName                    string `gorm:"column:username"`
	Tokens                      string `gorm:"column:tokens"`
	GeofenceName                string `gorm:"column:geofencename"`
	Attributes                  string `gorm:"column:attributes"`
	NotificationDoor            int    `gorm:"column:door"`
	NotificationPowerOn         int    `gorm:"column:ignitionOn"`
	NotificationPowerOff        int    `gorm:"column:ignitionOff"`
	NotificationDeviceOverspeed int    `gorm:"column:deviceOverspeed"`
	NotificationGeofenceEnter   int    `gorm:"column:geofenceEnter"`
	NotificationGeofenceExit    int    `gorm:"column:geofenceExit"`
	NotificationShock           int    `gorm:"column:shock"`
	NotificationPowerCut        int    `gorm:"column:powerCut"`
	NotificationLowBattery      int    `gorm:"column:lowBattery"`
	NotificationSos             int    `gorm:"column:sos"`
	Address                     string `gorm:"column:address"`
	Equivalent                  string
}

// esto existe para validar que el usuario tenga activo el tipo de notificacion y enviarla si es el caso
func StructToMap(e Events) map[string]int {
	result := make(map[string]int)
	result["door"] = e.NotificationDoor
	result["ignitionOn"] = e.NotificationPowerOn
	result["ignitionOff"] = e.NotificationPowerOff
	result["powerOn"] = e.NotificationPowerOn
	result["powerOff"] = e.NotificationPowerOff
	result["deviceOverspeed"] = e.NotificationDeviceOverspeed
	result["geofenceEnter"] = e.NotificationGeofenceEnter
	result["geofenceExit"] = e.NotificationGeofenceExit
	result["shock"] = e.NotificationShock
	result["powerCut"] = e.NotificationPowerCut
	result["lowBattery"] = e.NotificationLowBattery
	result["sos"] = e.NotificationSos
	return result
}

type MessageStatusResponse struct {
	AppName string
	Message string
	Token   string
}
type MessageStatus struct {
	Message *messaging.Message
	Event   Events
}

func FormatStatusMessage(event Events, success bool, aditionInfo string) MessageStatusResponse {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	statusMsg := "NO se envía PUSH"
	if success {
		statusMsg = "se envía PUSH"
	}
	return MessageStatusResponse{
		Message: fmt.Sprintf("%s %s %d %s a %s Dispositivo %s: Alerta Vehículo %s %s. %s", formattedTime, event.AppName, event.ID, statusMsg, event.UserName, event.Plate, event.Plate, event.Event, aditionInfo),
		AppName: event.AppName,
		Token:   event.Tokens,
	}
}
