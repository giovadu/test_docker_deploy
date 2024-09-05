package models

type UserInfo struct {
	Plate   string `gorm:"column:plate"`
	Name    string `gorm:"column:name"`
	Token   string `gorm:"column:token"`
	So      string `gorm:"column:so"`
	AppName string `gorm:"column:appname"`
}
