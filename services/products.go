package services

import (
	"myshipper/infrastructure"
	"myshipper/models"
)

func FetchProductsPage(page int, pageSize int) ([]models.Product, int, []int, error) {
	database := infrastructure.GetDb()
	var products []models.Product
	var count int
	tx := database.Begin()
	database.Model(&products).Count(&count)
	database.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
	tx.Model(&products).Preload("Tags").Preload("Categories").Preload("Images").Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
	commentCount := make([]int, len(products))
	for index, product := range products {
		commentCount[index] = tx.Model(&product).Association("Comments").Count()
	}
	err := tx.Commit().Error
	return products, count, commentCount, err
}

func FetchProductsDetails(condition interface{}, optional ...bool) models.Product {
	database := models.DB
	var product models.Product

	query := database.Where(condition).Preload("Tags").Preload("Categories").Preload("Images").Preload("Comments")
	query.First(&product)
	includeUserComment := false

	if len(optional) > 0 {
		includeUserComment = optional[0]
	}

	if includeUserComment {
		for i := 0; i < len(product.Comments); i++ {
			database.Model(&product.Comments[i]).Related(&product.Comments[i].User, "UserId")
		}

		var userIds = make([]uint, len(product.Comments))
		var users []models.User
		for i := 0; i < len(product.Comments); i++ {
			userIds[i] = product.Comments[i].UserId
		}
		database.Select("id, username").Where(userIds).Find(&users)

		for i := 0; i < len(product.Comments); i++ {
			user := users[i]
			comment := product.Comments[i]
			if comment.UserId == user.ID {
				product.Comments[i].User = user
			}
		}
	}
	return product
}

func FetchProductId(slug string) (uint, error) {
	productId := -1
	database := models.DB
	err := database.Model(&models.Product{}).Where(&models.Product{Slug: slug}).Select("id").Row().Scan(&productId)
	return uint(productId), err
}

func SetTags(product *models.Product, tags []string) error {
	database := models.DB
	var tagList []models.Tag
	for _, tag := range tags {
		var tagModel models.Tag
		err := database.FirstOrCreate(&tagModel, models.Tag{Name: tag}).Error
		if err != nil {
			return err
		}
		tagList = append(tagList, tagModel)
	}
	product.Tags = tagList
	return nil
}

func Update(product *models.Product, data interface{}) error {
	database := models.DB
	err := database.Model(product).Update(data).Error
	return err
}

func DeleteProduct(condition interface{}) error {
	database := models.DB
	err := database.Where(condition).Delete(models.Product{}).Error
	return err
}

func FetchProductsIdNameAndPrice(productIds []uint) (products []models.Product, err error) {
	database := models.DB
	err = database.Select([]string{"id", "name", "slug", "price"}).Find(&products, productIds).Error
	return products, err
}
