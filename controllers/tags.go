package controllers

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"myshipper/middlewares"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RegisterTagRoutes(router *gin.RouterGroup) {
	router.GET("", TagList)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("", CreateTag)
	}
}
