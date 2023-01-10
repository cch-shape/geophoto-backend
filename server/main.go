package main

import (
	"geophoto/backend/database"
	"geophoto/backend/middleware"
	"geophoto/backend/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Use(cors.New())

	database.Connect()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
