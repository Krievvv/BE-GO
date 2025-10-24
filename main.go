package main

import (
	"log"
	"os"
	"prak4/config"
	"prak4/database"
	"prak4/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()
	database.ConnectMongo()

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Alumni API")
	})

	route.SetupRoutes(app, database.DB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}