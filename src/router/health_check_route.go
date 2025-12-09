package router

import (
	"app/src/controller"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func HealthCheckRoutes(v1 fiber.Router, healthCheckService *service.HealthCheckService) {
	healthCheckController := controller.NewHealthCheckController(healthCheckService)

	healthCheck := v1.Group("/health-check")
	healthCheck.Get("/", healthCheckController.Check)
}
