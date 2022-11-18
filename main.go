package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"myshipper/controllers"
	"myshipper/infrastructure"
	"myshipper/middlewares"
	"myshipper/models"
	"myshipper/seeds"
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
		&models.ProductTag{},
		&models.Tag{},
		&models.Category{},
		&models.Product{},
		&models.UserRole{},
		&models.Role{},
		&models.User{},
	)
}
func migrate(database *gorm.DB) {
	database.AutoMigrate(&models.Address{})
	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.Comment{})
	database.AutoMigrate(&models.Order{})
	database.AutoMigrate(&models.OrderItem{})
	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.ProductCategory{})
	database.AutoMigrate(&models.Role{})
	database.AutoMigrate(&models.UserRole{})
	database.AutoMigrate(&models.Tag{})
	database.AutoMigrate(&models.ProductTag{})
	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.FileUpload{})
}

func addDbConstraints(database *gorm.DB) {
	database.Model(&models.Comment{}).AddForeignKey("product_id", "products(id)", "CASCADE", "CASCADE")
	database.Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	database.Model(&models.Order{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	database.Model(&models.Order{}).AddForeignKey("address_id", "addresses(id)", "CASCADE", "CASCADE")
	database.Model(&models.OrderItem{}).AddForeignKey("order_id", "orders(id)", "CASCADE", "CASCADE")
	database.Model(&models.OrderItem{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	database.Model(&models.Address{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	database.Model(&models.UserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	database.Model(&models.UserRole{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

	database.Table("products_tags").AddForeignKey("product_id", "products(id)", "CASCADE", "CASCADE")
	database.Table("products_tags").AddForeignKey("tag_id", "tags(id)", "CASCADE", "CASCADE")

	database.Model(models.ProductCategory{}).AddForeignKey("product_id", "products(id)", "CASCADE", "CASCADE")
	database.Model(models.ProductCategory{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")

	database.Model(&models.UserRole{}).AddIndex("user_roles__idx_user_id", "user_id")
	database.Table("products_tags").AddIndex("products_tags__idx_product_id", "product_id")
}

func create(database *gorm.DB) {
	drop(database)
	migrate(database)
	addDbConstraints(database)
	seeds.ImportProduct(database)
	os.Exit(0)
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
		} else if first == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if second == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if first != "" && second == "" {
			os.Exit(0)
		}
	}
	migrate(database)
	seeds.ImportProduct(database)
	goGonicEngine := gin.Default()

	goGonicEngine.Use(cors.Default())
	goGonicEngine.Use(middlewares.Benchmark())
	goGonicEngine.Use(middlewares.UserLoaderMiddleware())

	goGonicEngine.Static("/static", "./static")

	apiRouteGroup := goGonicEngine.Group("/api")

	controllers.RegisterUserRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterProductRoutes(apiRouteGroup.Group("/products"))
	controllers.RegisterCommentRoutes(apiRouteGroup.Group("/"))
	controllers.RegisterPageRoutes(apiRouteGroup.Group("/"))
	controllers.RegisterAddressesRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterTagRoutes(apiRouteGroup.Group("/tags"))
	controllers.RegisterCategoryRoutes(apiRouteGroup.Group("/categories"))
	controllers.RegisterOrderRoutes(apiRouteGroup.Group("/orders"))

	goGonicEngine.Run(":5052")

}
