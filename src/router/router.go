package router

import (
	"app/src/config"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	container := dig.New()

	container.Provide(func() *gorm.DB { return db })
	container.Provide(validation.Validator)

	container.Provide(service.NewHealthCheckService)
	container.Provide(service.NewEmailService)
	container.Provide(service.NewUserService)
	container.Provide(service.NewTokenService)
	container.Provide(service.NewAuthService)

	// Invoke route setup with auto-injected dependencies
	err := container.Invoke(func(
		healthCheckService *service.HealthCheckService,
		emailService *service.EmailService,
		userService *service.UserService,
		tokenService *service.TokenService,
		authService *service.AuthService,
	) {
		v1 := app.Group("/v1")

		HealthCheckRoutes(v1, healthCheckService)
		AuthRoutes(v1, authService, userService, tokenService, emailService)
		UserRoutes(v1, userService, tokenService)

		if !config.IsProd {
			DocsRoutes(v1)
		}
	})

	if err != nil {
		panic(err)
	}
}
