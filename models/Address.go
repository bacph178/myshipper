package models

import "github.com/jinzhu/gorm"

type Address struct {
	gorm.Model
	StreetAddress string  `gorm:"not null"`
	City          string  `gorm:"not null"`
	Country       string  `gorm:"not null"`
	ZipCode       string  `gorm:"not null"`
	FirstName     string  `gorm:"not null"`
	LastName      string  `gorm:"not null"`
	User          User    `gorm:"association_foreignkey:UserId:"`
	UserId        uint    `gorm:"default:null"`
	Orders        []Order `gorm:"foreignKey:AddressId"`
}
