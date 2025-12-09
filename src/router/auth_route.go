package router

import (
	"app/src/config"
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(
	v1 fiber.Router, authService *service.AuthService, userService *service.UserService,
	tokenService *service.TokenService, emailService *service.EmailService,
) {
	authController := controller.NewAuthController(authService, userService, tokenService, emailService)
	config.GoogleConfig()

	auth := v1.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Post("/refresh-tokens", authController.RefreshTokens)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/send-verification-email", middleware.Auth(userService), authController.SendVerificationEmail)
	auth.Post("/verify-email", authController.VerifyEmail)
	auth.Get("/google", authController.GoogleLogin)
	auth.Get("/google-callback", authController.GoogleCallback)
}
