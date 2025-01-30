package routes

import (
	"github.com/chamanbravo/upstat/internal/controllers/rest"
	"github.com/chamanbravo/upstat/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// @Group Monitors
func MonitorRoutes(app *fiber.App, h *controllers.Handler) {
	route := app.Group("/api/monitors", middleware.Protected)

	route.Post("", h.CreateMonitor)
	route.Get("", h.MonitorsList)
	route.Get("/:id", h.MonitorInfo)
	route.Patch("/:id", h.UpdateMonitor)
	route.Delete("/:id", h.DeleteMonitor)
	route.Patch(":id/pause", h.PauseMonitor)
	route.Patch(":id/resume", h.ResumeMonitor)
	route.Get("/:id/summary", h.MonitorSummary)
	route.Get("/:id/heartbeat", h.RetrieveHeartbeat)
	route.Get("/:id/cert-exp-countdown", h.CertificateExpiryCountDown)
	route.Get("/:id/notifications", h.NotificationChannelListOfMonitor)
	route.Get("/:id/status-pages", h.StatusPagesListOfMonitor)
}
