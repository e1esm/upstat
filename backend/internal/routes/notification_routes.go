package routes

import (
	"github.com/chamanbravo/upstat/internal/controllers/rest"
	"github.com/chamanbravo/upstat/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// @Group Notifications
func NotificationRoutes(app *fiber.App, h *controllers.Handler) {
	route := app.Group("/api/notifications", middleware.Protected)

	route.Post("", h.CreateNotification)
	route.Get("", h.ListNotificationsChannel)
	route.Delete("/:id", h.DeleteNotificationChannel)
	route.Patch("/:id", h.UpdateNotificationChannel)
	route.Get("/:id", h.NotificationChannelInfo)
}
