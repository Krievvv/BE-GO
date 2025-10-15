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
	// Memuat variabel lingkungan
	config.LoadEnv()

	// Terhubung ke database
	database.ConnectDB()

	// Membuat instance aplikasi Fiber baru
	app := fiber.New()
	app.Use(logger.New())

	// Tambahkan route untuk root agar tidak 404
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Alumni API")
	})

	// Mengatur rute (routes)
	route.SetupRoutes(app, database.DB)

	// Memulai server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}