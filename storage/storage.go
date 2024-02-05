package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zykunov/courseGoFirst/vkApiBot/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// func NewConnection(c *ConnectionDB) (db *gorm.DB) {
func init() {

	if err := godotenv.Load("configs/.env"); err != nil {
		log.Print("No .env file found")
	}

	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("Dbname")
	sslmode := os.Getenv("sslmode")

	var dbString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	log.Println("try connect to DB")

	connection, err := gorm.Open(postgres.Open(dbString), &gorm.Config{})
	if err != nil {
		log.Fatalf("can't connect to DB", err)
	}

	db = connection
	db.Debug().AutoMigrate(&models.Group{})
	db.Debug().AutoMigrate(&models.User{})

}

func GetDB() *gorm.DB {
	return db
}
