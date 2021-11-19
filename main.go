package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/vikas/config"
	"github.com/vikas/routes"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())
	routes.Route(app)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading .env file")
	}

	config.ConnectDB()

	if err = app.Listen(":8080"); err != nil {
		log.Fatal("can not listen")
	}

}
