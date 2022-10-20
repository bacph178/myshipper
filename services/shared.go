package services

import "myshipper/models"

func CreateOne(data interface{}) error {
	database := models.DB
	err := database.Create(data).Error
	return err
}
