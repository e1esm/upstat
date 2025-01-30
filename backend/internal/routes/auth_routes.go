package routes

import (
	"github.com/chamanbravo/upstat/internal/controllers/rest"

	"github.com/gofiber/fiber/v2"
)

// @Group Auth
func AuthRoutes(app *fiber.App, h *controllers.Handler) {
	route := app.Group("/api/auth")

	route.Post("/signup", h.SignUp)
	route.Post("/signin", h.SignIn)
	route.Post("/signout", h.SignOut)
	route.Post("/refresh-token", h.RefreshToken)
}
