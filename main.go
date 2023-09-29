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
	db := config.DBInit()

	app := fiber.New()
	router.StartApp(app, db)

	app.Listen(":"+config.GetEnv("PORT"))

}
