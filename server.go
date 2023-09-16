package main

import (
	"github.com/Abhijeet6387/Blog/database"
	"github.com/Abhijeet6387/Blog/routes"
	"github.com/gofiber/fiber/v2"
)
func main(){

	database.Connect()
	app := fiber.New()
	routes.Setup(app)
	app.Listen(":5000")
}