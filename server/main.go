package main

import (
	"geophoto/backend/database"
	"geophoto/backend/middleware"
	"geophoto/backend/router"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    20 * 1024 * 1024,
	})
	app.Use(cors.New())

	database.Connect()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
