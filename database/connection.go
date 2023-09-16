package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Abhijeet6387/Blog/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(){
	envErr := godotenv.Load(".env")
	if envErr != nil{
		log.Fatal("Unable to load .env file!")
	}
	DBname := os.Getenv("DB_NAME")
	Password := os.Getenv("PASSWORD")

	dsn := fmt.Sprintf("root:%s@/%s", Password, DBname)
	// fmt.Printf("%T : %s", dsn, dsn)

	connection , dbErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if dbErr != nil{
		log.Fatal("Could'nt connect with database!")
	}

	// we need database in controllers, declare a global variable *gorm.DB
	DB = connection  

	// we need to migrate the model to create table in database 
	connection.AutoMigrate(&models.User{})			
}