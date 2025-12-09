package controller

import (
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckController struct {
	healthCheckService *service.HealthCheckService
}

func NewHealthCheckController(healthCheckService *service.HealthCheckService) *HealthCheckController {
	return &HealthCheckController{
		healthCheckService,
	}
}

func (h *HealthCheckController) addServiceStatus(
	serviceList *[]response.HealthCheck, name string, isUp bool, message *string,
) {
	status := "Up"

	if !isUp {
		status = "Down"
	}

	*serviceList = append(*serviceList, response.HealthCheck{
		Name:    name,
		Status:  status,
		IsUp:    isUp,
		Message: message,
	})
}

// @Tags Health
// @Summary Health Check
// @Description Check the status of services and database connections
// @Accept json
// @Produce json
// @Success 200 {object} example.HealthCheckResponse
// @Failure 500 {object} example.HealthCheckResponseError
// @Router /health-check [get]
func (h *HealthCheckController) Check(c *fiber.Ctx) error {
	isHealthy := true
	var serviceList []response.HealthCheck

	// Helper to run a health check, update status and add to list
	checkStatus := func(name string, fn func() error) {
		err := fn()
		if err != nil {
			isHealthy = false
			errMsg := err.Error()
			h.addServiceStatus(&serviceList, name, false, &errMsg)
		} else {
			h.addServiceStatus(&serviceList, name, true, nil)
		}
	}

	checkStatus("Postgre", h.healthCheckService.GormCheck)
	checkStatus("Memory", h.healthCheckService.MemoryHeapCheck)

	statusCode := fiber.StatusOK
	status := "success"
	if !isHealthy {
		statusCode = fiber.StatusInternalServerError
		status = "error"
	}

	return c.Status(statusCode).JSON(response.HealthCheckResponse{
		Status:    status,
		Message:   "Health check completed",
		Code:      statusCode,
		IsHealthy: isHealthy,
		Result:    serviceList,
	})
}
