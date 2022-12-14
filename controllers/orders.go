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

func RegisterOrderRoutes(router *gin.RouterGroup) {
	router.POST("", CreateOrder)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("", ListOrders)
		router.GET("/:id", ShowOrder)
	}
}

func ListOrders(c *gin.Context) {
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	userId := c.MustGet("currentUserId").(uint)
	orders, totalCommentCount, err := services.FetchOrdersPage(userId, page, pageSize)
	c.JSON(http.StatusOK, dtos.CreateOrderPagedResponse(c.Request, orders, page, pageSize, totalCommentCount, false, false))
}

func ShowOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	user := c.MustGet("currentUser").(models.User)
	order, err := services.FetchOrderDetails(uint(orderId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	if order.UserId == user.ID || user.IsAdmin() {
		c.JSON(http.StatusOK, dtos.CreateOrderDetailsDto(&order))
	} else {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you can not view this order"))
		return
	}
}

func CreateOrder(c *gin.Context) {
	var orderRequest dtos.CreateOrderRequestDto
	if err := c.ShouldBind(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	userObj, userLoggedIn := c.Get("currentUser")
	var user models.User
	if userLoggedIn {
		user = (userObj).(models.User)
	}
	var address models.Address
	if orderRequest.AddressId != 0 && userLoggedIn {
		address = services.FetchAddress(orderRequest.AddressId)
		if address.UserId != user.ID {
			c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("permission denied"))
			return
		}
	} else if orderRequest.AddressId == 0 {
		address = models.Address{
			FirstName:     orderRequest.FirstName,
			LastName:      orderRequest.LastName,
			City:          orderRequest.City,
			Country:       orderRequest.Country,
			StreetAddress: orderRequest.StreetAddress,
			ZipCode:       orderRequest.ZipCode,
		}
		if userLoggedIn {
			address.UserId = user.ID
		}
		err := services.CreateOne(&address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	} else {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Operation not supported, what are you trying to do?"))
		return
	}
	order := models.Order{
		TrackingNumber: randomString(16),
		OrderStatus:    0,
		Address:        address,
		AddressId:      address.ID,
	}

	if userLoggedIn {
		order.UserId = user.ID
		order.User = user
	}
	var productIds = make([]uint, len(orderRequest.CartItems))
	for i := 0; i < len(orderRequest.CartItems); i++ {
		productIds[i] = orderRequest.CartItems[i].Id
	}
	products, err := services.FetchProductsIdNameAndPrice(productIds)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	if len(products) != len(orderRequest.CartItems) {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateErrorDtoWithMessage("make sure all products are still available"))
		return
	}
	orderItems := make([]models.OrderItem, len(products))
	for i := 0; i < len(products); i++ {
		orderItems[i] = models.OrderItem{
			ProductID:   products[i].ID,
			ProductName: products[i].Name,
			Slug:        products[i].Slug,
			Quantity:    orderRequest.CartItems[i].Quantity,
		}
	}
	order.OrderItems = orderItems
	err = services.CreateOne(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, dtos.CreateOrderCreatedDto(&order))
}
