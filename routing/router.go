package routing

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	_ "go-template/docs"
	"go-template/models"
	"go-template/service/app"
)

func Setup(f *fiber.App) {
	f.Get("/swagger/*", swagger.HandlerDefault)
	appApi := f.Group("/api")
	appApi.Get("/user/findUserByNameAndPwd", app.FindUserByNameAndPwd)
	appApi.Get("/user/getUserList", app.GetUserList)
	appApi.Get("/user/createUser", app.CreateUser)
	appApi.Get("/user/deleteUser", app.DeleteUser)
	appApi.Post("/user/updateUser", app.UpdateUser)
	appApi.Get("/user/sendMsg", app.SendMsg)
	appApi.Get("/user/chat", models.Chat)
}
