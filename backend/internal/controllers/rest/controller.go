package controllers

import "github.com/chamanbravo/upstat/internal/app"

type Handler struct {
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{
		app: app,
	}
}
