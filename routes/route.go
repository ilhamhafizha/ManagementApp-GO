package routes

func (app *fiber.App,uc *controller.UserController) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found.")
	}
	
	app.Post("v1/api/register", uc.Register)
}