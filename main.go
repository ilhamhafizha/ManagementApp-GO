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

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	app := fiber.New()

	app.All("*", func(c *fiber.Ctx) error {
		log.Println("METHOD:", c.Method(), "PATH:", c.Path())
		return c.Next()
	})

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boarController := controllers.NewBoardController(boardService)

	routes.Setup(app, userController, boarController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port : " + port)
	log.Fatal(app.Listen(":" + port))

}
