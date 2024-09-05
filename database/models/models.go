package models

// Events representa la estructura de un evento en la base de datos.
type Events struct {
	ID           int     `gorm:"column:id" json:"id"`
	ServerTime   string  `gorm:"column:servertime" json:"servertime"`
	Type         string  `gorm:"column:type" json:"type"`
	Event        string  `gorm:"column:event" json:"event"`
	Plate        string  `gorm:"column:plate" json:"plate"`
	GeofenceName *string `gorm:"column:geofencename" json:"geofencename"`
}
type VerificationMessage struct {
	AppName   string `json:"app_name"`   // Nombre de la aplicación
	Message   string `json:"message"`    // Mensaje de verificación
	EventID   int    `json:"event_id"`   // ID del evento asociado
	CreatedAt string `json:"created_at"` // Fecha de creación del mensaje
}
