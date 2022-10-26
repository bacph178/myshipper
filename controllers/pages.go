package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"myshipper/dtos"
	"myshipper/services"
	"net/http"
)

func RegisterPageRoutes(router *gin.RouterGroup) {
	router.GET("", Home)
	router.GET("/home", Home)
}

func Home(c *gin.Context) {
	tags, tagErr := services.FetchAllTags()
	categories, CatErr := services.FetchAllCategories()
	if tagErr != nil || CatErr != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("comments", errors.New("Something went wrong")))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateHomeResponse(tags, categories))
}
