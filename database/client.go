package database

import (
	"cinetodoapi/model"
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
	Migrate()
	log.Println("Connected to Database!")
}

func Migrate() {
	Instance.AutoMigrate(&model.User{}, &model.Movie{}, &model.UserMovie{})
	log.Println("Database Migration Completed!")
}
