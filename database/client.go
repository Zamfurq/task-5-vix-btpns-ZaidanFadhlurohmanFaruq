package database

import (
	"github/golang-api/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

func Connect(connectionString string) {
	Instance, dbError = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to database!")
}
func Migrate() {
	Instance.AutoMigrate(&model.User{})
	Instance.AutoMigrate(&model.Photo{})
	log.Println("Database Migration Completed!")
}
