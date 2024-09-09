package models

// Events representa la estructura de un evento en la base de datos.
type Events struct {
	ID           int    `gorm:"column:id" json:"id"`
	Type         string `gorm:"column:type" json:"type"`
	Event        string `gorm:"column:event" json:"event"`
	AppName      string `gorm:"column:appname"`
	Plate        string `gorm:"column:name"`
	UserName     string `gorm:"column:username"`
	Tokens       string `gorm:"column:tokens"`
	GeofenceName string `gorm:"column:geofencename"`
}
