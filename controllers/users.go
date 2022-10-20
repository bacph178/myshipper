package controllers

import (
	"github.com/gin-gonic/gin"
	"myshipper/dtos"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

func UsersRegistration(c *gin.Context) {
	var json dtos.Re
}
