package router

import (
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(v1 fiber.Router, userService *service.UserService, tokenService *service.TokenService) {
	userController := controller.NewUserController(userService, tokenService)

	user := v1.Group("/users")

	user.Get("/", middleware.Auth(userService, "getUsers"), userController.GetUsers)
	user.Post("/", middleware.Auth(userService, "manageUsers"), userController.CreateUser)
	user.Get("/:userId", middleware.Auth(userService, "getUsers"), userController.GetUserByID)
	user.Patch("/:userId", middleware.Auth(userService, "manageUsers"), userController.UpdateUser)
	user.Delete("/:userId", middleware.Auth(userService, "manageUsers"), userController.DeleteUser)
}
