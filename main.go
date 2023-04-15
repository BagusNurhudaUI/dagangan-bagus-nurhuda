package main

import (
	"fmt"

	"github.com/BagusNurhudaUI/dagangan-bagus-nurhuda/config"
	"github.com/BagusNurhudaUI/dagangan-bagus-nurhuda/router"

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
