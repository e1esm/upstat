package main

import (
	"log"

	_ "github.com/chamanbravo/upstat/docs"
	appLayer "github.com/chamanbravo/upstat/internal/app"
	controllers "github.com/chamanbravo/upstat/internal/controllers/rest"
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

	repo := repository.New(db)

	monitor := pkg.New(repo)

	aLayer := appLayer.New(repo, monitor)

	h := controllers.New(aLayer)

	monitor.StartGoroutineSetup()

	routes.AuthRoutes(app, h)
	routes.SwaggerRoute(app)
	routes.MonitorRoutes(app, h)
	routes.UserRoutes(app, h)
	routes.NotificationRoutes(app, h)
	routes.StatusPagesRoutes(app, h)

	log.Fatal(app.Listen(":8000"))
}
