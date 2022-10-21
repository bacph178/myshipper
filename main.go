package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"myshipper/controllers"
	"myshipper/infrastructure"
	"myshipper/middlewares"
	"myshipper/models"
	"os"
)

func drop(database *gorm.DB) {
	database.DropTableIfExists(
		&models.FileUpload{},
		&models.Comment{},
		&models.OrderItem{},
		&models.Order{},
		&models.Address{},
		&models.ProductCategory{},
		&models.Product{},
		&models.UserRole{},
		&models.Role{},
		&models.User{},
	)
}

func create(database *gorm.DB) {
	drop(database)
}

func main() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	println(os.Getenv("DB_DIALECT"))
	database := infrastructure.OpenDbConnection()
	defer database.Close()
	args := os.Args
	if len(args) > 1 {
		first := args[1]
		second := ""
		if len(args) > 2 {
			second = args[2]
		}
		if first == "create" {
			create(database)
		}
	}

	//models.ConnectDataBase()

	goGonicEngine := gin.Default()

	//public := goGonicEngine.Group("/api")

	//public.POST("/register", controllers.Register)
	//public.POST("/login", controllers.Login)

	//protected := goGonicEngine.Group("/api/admin")
	//protected.Use(middlewares.JwtAuthMiddleware())
	//protected.GET("/user", controllers.CurrentUser)

	goGonicEngine.Use(middlewares.U)
	apiRouteGroup := goGonicEngine.Group("/api")
	controllers.RegisterProductRoutes(apiRouteGroup.Group("/products"))

	goGonicEngine.Run(":5052")

}
