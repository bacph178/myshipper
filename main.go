package main

import (
	"myshipper/controllers"
	"myshipper/middlewares"
	"myshipper/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDataBase()

	goGonicEngine := gin.Default()

	public := goGonicEngine.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := goGonicEngine.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	//protected.GET("/user", controllers.CurrentUser)

	apiRouteGroup := goGonicEngine.Group("/api")
	controllers.RegisterProductRoutes(apiRouteGroup.Group("/products"))

	goGonicEngine.Run(":5052")

}
