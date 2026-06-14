package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/ilhamhafizha/ManagementApp-GO/config"
	"github.com/ilhamhafizha/ManagementApp-GO/controllers"
	"github.com/ilhamhafizha/ManagementApp-GO/utils"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App,
	uc *controllers.UserController,
	bc *controllers.BoardController,
	lc *controllers.ListController,
	cc *controllers.CardController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found.")
	}

	app.Post("/api/v1/register", uc.Register)
	app.Post("/api/v1/login", uc.Login)

	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error Unauthorized", err.Error())
		},
	}))

	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.GetUserPagination)
	userGroup.Get("/:id", uc.GetUser)
	userGroup.Put("/:id", uc.UpdateUser)
	userGroup.Delete("/:id", uc.DeleteUser)

	boardGroup := api.Group("/boards")
	boardGroup.Post("/", bc.CreateBoard)
	boardGroup.Put("/:id", bc.UpdateBoard)
	boardGroup.Post("/:id/members", bc.AddBoardMember)
	boardGroup.Delete("/:id/members", bc.RemoveBoardMember)
	boardGroup.Get("/my", bc.GetMyBoardPaginate)
	boardGroup.Get("/:board_id/lists", lc.GetListOnBoard)
	boardGroup.Put("/:board_id/position", lc.UpdateListPosition)

	listGroup := api.Group("/lists")
	listGroup.Post("/", lc.CreateList)
	listGroup.Put("/:id", lc.UpdateList)
	listGroup.Delete("/:id", lc.DeleteList)

	listGroup.Get("/:list_id/cards", cc.GetListCard)

	//card
	cardGroup := api.Group("/cards")
	cardGroup.Post("/", cc.CreateCard)
	cardGroup.Put("/:id", cc.UpdateCard)
	cardGroup.Delete("/:id", cc.DeleteCard)
	cardGroup.Get("/:id", cc.GetCardDetail)
}
