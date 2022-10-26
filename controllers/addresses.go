package controllers

import (
	"github.com/gin-gonic/gin"
	"myshipper/dtos"
	"myshipper/middlewares"
	"myshipper/models"
	"myshipper/services"
	"net/http"
	"strconv"
)

func RegisterAddressesRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("/addresses", ListAddresses)
		router.POST("/addresses", CreateAddress)
	}

}

func ListAddresses(c *gin.Context) {
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	pageSize, pageSizeErr := strconv.Atoi(pageSizeStr)
	if pageSizeErr != nil {
		pageSize = 5
	}

	page, pageErr := strconv.Atoi(pageStr)
	if pageErr != nil {
		page = 1
	}

	userId := c.MustGet("currentUserId").(uint)
	includeUser := false
	addresses, totalCommentCount := services.FetchAddressesPage(userId, page, pageSize, includeUser)
	c.JSON(http.StatusOK, dtos.CreateAddressPagedResponse(c.Request, addresses, page, pageSize, totalCommentCount, includeUser))
}

func CreateAddress(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	var json dtos.CreateAddress
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}
	firstName := json.FirstName
	lastName := json.LastName
	if firstName == "" {
		firstName = user.FirstName
	}
	if lastName == "" {
		lastName = user.LastName
	}
	address := models.Address{
		FirstName:     firstName,
		LastName:      lastName,
		Country:       json.Country,
		City:          json.City,
		StreetAddress: json.StreetAddress,
		ZipCode:       json.ZipCode,
		User:          user,
		UserId:        user.ID,
	}
	if err := services.SaveOne(&address); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database_error", err))
		return
	}
	c.JSON(http.StatusOK, dtos.GetAddressCreatedDto(&address, false))
}
