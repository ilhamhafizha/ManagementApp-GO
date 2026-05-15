package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ilhamhafizha/ManagementApp-GO/config"
	"github.com/ilhamhafizha/ManagementApp-GO/controllers"
	"github.com/ilhamhafizha/ManagementApp-GO/database/seed"
	"github.com/ilhamhafizha/ManagementApp-GO/repositories"
	"github.com/ilhamhafizha/ManagementApp-GO/routes"
	"github.com/ilhamhafizha/ManagementApp-GO/services"
)

func main(){
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	routes.Setup(app, userController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port : " + port)
	log.Fatal(app.Listen(":" + port))

}