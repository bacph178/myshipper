package infrastructure

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"path"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func OpenDbConnection() *gorm.DB {
	dialect := os.Getenv("DB_DIALECT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	//port := os.Getenv("DB_PORT")
	var db *gorm.DB
	var err error
	databaseUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable ", host, username, password, dbName)
	//databaseUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
	println(databaseUrl)
	db, err = gorm.Open(dialect, databaseUrl)
	if err != nil {
		fmt.Println("db err: ", err)
		os.Exit(-1)
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	DB = db
	return DB
}

func RemoveDB(db *gorm.DB) error {
	db.Close()
	err := os.Remove(path.Join(".", "app.db"))
	return err
}

func GetDb() *gorm.DB {
	return DB
}
