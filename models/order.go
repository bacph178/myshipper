package models

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.DB
	OrderStatus     int `gorm:"default:0"`
	TrackingNumber  string
	OrderItems      []OrderItem `gorm:"foreignKey:OrderId"`
	Address         Address     `gorm:"association_foreignkey:AddressId:"`
	AddressId       uint
	User            User `gorm:"foreignKey:UserId:"`
	UserId          uint `gorm:"default:null"`
	OrderItemsCount int  `gorm:"-"`
}
