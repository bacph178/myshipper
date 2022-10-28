package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"myshipper/dtos"
	"myshipper/infrastructure"
	"myshipper/middlewares"
	"myshipper/models"
	"myshipper/services"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterCategoryRoutes(router *gin.RouterGroup) {
	router.GET("", CategoryList)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("", CreateCategory)
	}
}

func CategoryList(c *gin.Context) {
	categories, err := services.FetchAllCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("fetch_error", err))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateCategoryListMapDto(categories))
}

func CreateCategory(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you must be admin"))
		return
	}
	name := c.PostForm("name")
	description := c.PostForm("description")
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form error: %s", err))
		return
	}
	files := form.File["images[]"]
	var categoryImages = make([]models.FileUpload, len(files))
	for index, file := range files {
		filename := randomString(16) + ".png"
		dirpath := filepath.Join(".", "static", "images", "categories")
		imagePath := filepath.Join(dirpath, filename)
		if _, err = os.Stat(dirpath); os.IsNotExist(err) {
			err = os.MkdirAll(dirpath, os.ModeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("io_error", err))
				return
			}
		}

		outputFile, err := os.Create(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		inputFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusOK, dtos.CreateDetailedErrorDto("io_error", err))
		}
		defer inputFile.Close()

		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		fileSize := (uint)(file.Size)
		categoryImages[index] = models.FileUpload{Filename: file.Filename, FilePath: string(filepath.Separator) + imagePath, FileSize: fileSize}
	}
	database := infrastructure.GetDb()
	category := models.Category{Name: name, Description: description, Images: categoryImages}
	err = database.Create(&category).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateCategoryCreatedDto(category))
}
