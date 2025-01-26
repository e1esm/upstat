package routes

import (
	"github.com/chamanbravo/upstat/internal/controllers/rest"
	"github.com/chamanbravo/upstat/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// @Group Users
func UserRoutes(app *fiber.App, h *controllers.Handler) {
	route := app.Group("/api/users")

	route.Get("/setup", h.Setup)
	route.Post("/update-password", middleware.Protected, h.UpdatePassword)
	route.Patch("/me", middleware.Protected, h.UpdateAccount)
}
