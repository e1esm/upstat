package routes

import (
	"github.com/chamanbravo/upstat/internal/controllers/rest"
	"github.com/chamanbravo/upstat/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// @Group StatusPages
func StatusPagesRoutes(app *fiber.App, h *controllers.Handler) {
	route := app.Group("/api/status-pages", middleware.Protected)

	route.Post("", h.CreateStatusPage)
	route.Get("", h.ListStatusPages)
	route.Delete("/:id", h.DeleteStatusPage)
	route.Patch("/:id", h.UpdateStatusPage)
	route.Get("/:id", h.StatusPageInfo)
	route.Get("/:slug/summary", h.StatusSummary)
}
