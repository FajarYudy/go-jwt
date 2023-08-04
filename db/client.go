package db

import (
	"go-jwt/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Gorm *gorm.DB
var dbError error

func Connect(connectionString string){
	Gorm, dbError = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database!")

}

func Migrate(){
	Gorm.AutoMigrate(&models.User{})
	log.Println("Database Migration Completed!")
}