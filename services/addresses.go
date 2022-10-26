package services

import (
	"myshipper/infrastructure"
	"myshipper/models"
)

func FetchAddressesPage(userId uint, page, pageSize int, includeUser bool) ([]models.Address, int) {
	var addresses []models.Address
	var totalAddressesCount int
	database := infrastructure.GetDb()
	database.Model(&models.Address{}).Where(&models.Address{UserId: uint(userId)}).Count(&totalAddressesCount)
	database.Where(&models.Address{UserId: uint(userId)}).
		Offset((page - 1) * pageSize).Limit(pageSize).Preloads("User").Find(&addresses)
	if includeUser {
		var userIds = make([]uint, len(addresses))
		var users []models.User
		for i := 0; i < len(addresses); i++ {
			userIds[i] = addresses[i].UserId
		}
		database.Select([]string{"id", "username"}).Where(userIds).Find(&users)
		for i := 0; i < len(addresses); i++ {
			address := addresses[i]
			for j := 0; j < len(users); j++ {
				user := users[j]
				if address.UserId == user.ID {
					addresses[i].User = users[j]
				}
			}
		}
	}
	return addresses, totalAddressesCount
}

func FetchAddress(addressId uint) (address models.Address) {
	database := infrastructure.GetDb()
	database.First(&address, addressId)
	return address
}

func FetchIdsFromAddress(addressId uint) (address models.Address) {
	database := infrastructure.GetDb()
	database.Select("id, user_id").First(&address, addressId)
	return address
}
