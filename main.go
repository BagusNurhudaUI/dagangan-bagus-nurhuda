package main

import (
	"dagangan/config"
	"dagangan/router"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {

	fmt.Println("Starting...")

	//first initialize database
	config.DBInit()

	app := fiber.New()
	router.StartApp(app)

	app.Listen(config.GetEnv("PORT"))

}
