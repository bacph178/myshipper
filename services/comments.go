package services

import (
	"myshipper/infrastructure"
	"myshipper/models"
)

func FetchCommentsPage(productId, page int, pageSize int) ([]models.Comment, int) {
	var comments []models.Comment
	var totalCommentCount int
	database := infrastructure.GetDb()
	database.Model(&comments).Where(&models.Comment{ProductId: uint(productId)}).Count(&totalCommentCount)
	database.Where(&models.Comment{ProductId: uint(productId)}).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("User").
		Find(&comments)
	var userIds = make([]uint, len(comments))
	var users []models.User
	for i := 0; i < len(comments); i++ {
		userIds[i] = comments[i].UserId
	}
	database.Select("id, username").Where(userIds).Find(&users)
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		for j := 0; j < len(users); j++ {
			user := users[j]
			if comment.UserId == user.ID {
				comments[i].User = users[j]
			}
		}
	}
	return comments, totalCommentCount
}

func FetchCommentById(id int, includes ...bool) models.Comment {
	includeUser := false
	if len(includes) > 0 {
		includeUser = includes[0]
	}
	includeProduct := false
	if len(includes) > 1 {
		includeProduct = includes[1]
	}
	database := infrastructure.GetDb()
	var comment models.Comment
	if includeProduct && includeUser {
		database.Preload("User").Preload("Product").Find(&comment, id)
	} else if includeUser {
		database.Preload("User").Find(&comment, id)
	} else if includeProduct {
		database.Preload("Product").Find(&comment, id)
	} else {
		database.Find(&comment, id)
	}
	return comment
}

func DeleteComment(condition interface{}) error {
	database := infrastructure.GetDb()
	err := database.Where(condition).Delete(models.Comment{}).Error
	return err
}
