package routes

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamhafizha/ManagementApp-GO/controllers"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controllers.UserController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found.")
	}
	
	app.Post("v1/api/register", uc.Register)
} 