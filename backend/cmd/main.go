package main

import (
	"log"

	_ "github.com/chamanbravo/upstat/docs"
	"github.com/chamanbravo/upstat/internal/database"
	"github.com/chamanbravo/upstat/internal/repository"
	"github.com/chamanbravo/upstat/pkg"

	"github.com/chamanbravo/upstat/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

// @title Upstat API
// @version 1.0
// @description This is an auto-generated API Docs for Upstat API.
// @contact.email chamanpro9@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app := fiber.New()
	app.Use(logger.New())

	db, err := database.DBConnect()
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	_ = repository.New(db)

	pkg.StartGoroutineSetup()

	routes.AuthRoutes(app)
	routes.SwaggerRoute(app)
	routes.MonitorRoutes(app)
	routes.UserRoutes(app)
	routes.NotificationRoutes(app)
	routes.StatusPagesRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
