package main

import (
	"fmt"
	"log"
	"os"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT environment variable is not set
	}

	log.Printf("Server starting on :%s", port)
	app.Listen(":"+config.GetEnv("PORT"))

}
