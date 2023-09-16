package routes

import (
	"github.com/Abhijeet6387/Blog/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	
	app.Get("/", controllers.Home)
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/getuserdetails", controllers.GetUser)
	app.Post("/api/logout", controllers.Logout)

}