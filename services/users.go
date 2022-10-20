package services

import "myshipper/models"

func FindOneUser(condition interface{}) (models.User, error) {
	database := models.DB
	var user models.User
	err := database.Where(condition).Preload("Roles").First(&user).Error
	return user, err
}

func UpdateUser(user models.User, data interface{}) error {
	database := models.DB
	err := database.Model(user).Update(data).Error
	return err
}
