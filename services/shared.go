package services

import (
	"myshipper/infrastructure"
)

func CreateOne(data interface{}) error {
	database := infrastructure.GetDb()
	err := database.Create(data).Error
	return err
}
func SaveOne(data interface{}) error {
	database := infrastructure.GetDb()
	err := database.Save(data).Error
	return err
}
