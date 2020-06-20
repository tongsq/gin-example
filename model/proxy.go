package model

//import (
//	"github.com/jinzhu/gorm"
//)
type Proxy struct {
	Id         int    `gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY"`
	Host       string `gorm:"column:host"`
	Port       string
	Status     int8   `gorm:"column:status"`
	CreateTime int    `gorm:"column:create_time"`
	UpdateTime int    `gorm:"column:update_time"`
	ActiveTime int    `gorm:"column:active_time"`
	Country    string `gorm:"column:country"`
	Region     string `gorm:"column:region"`
	City       string `gorm:"column:city"`
	Isp        string `gorm:"column:isp"`
}

func (Proxy) TableName() string {
	return "proxy"
}
