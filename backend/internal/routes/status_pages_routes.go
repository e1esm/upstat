package routes

import (
	"github.com/chamanbravo/upstat/internal/controllers"
	"github.com/chamanbravo/upstat/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// @Group StatusPages
func StatusPagesRoutes(app *fiber.App) {
	route := app.Group("/api/status-pages", middleware.Protected)

	route.Post("", controllers.CreateStatusPage)
	route.Get("", controllers.ListStatusPages)
	route.Delete("/:id", controllers.DeleteStatusPage)
	route.Patch("/:id", controllers.UpdateStatusPage)
	route.Get("/:id", controllers.StatusPageInfo)
	route.Get("/:slug/summary", controllers.StatusSummary)
}
