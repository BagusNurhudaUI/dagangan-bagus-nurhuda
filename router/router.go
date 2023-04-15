package router

import (
	"github.com/BagusNurhudaUI/dagangan-bagus-nurhuda/controllers"

	"github.com/gofiber/fiber/v2"
)

func StartApp(app *fiber.App) {
	app.Get("/", GetIndex)
	products := app.Group("/products")
	products.Get("/", controllers.Paginate)                   // get all products with pagination
	products.Get("/all", controllers.GetProduct)              // get all products without pagination
	products.Get("/:productId", controllers.GetProductById)   // get all products by productId
	products.Post("/", controllers.PostProduct)               // create a new product
	products.Put("/:productId", controllers.UpdateProduct)    // update a product with existing data
	products.Delete("/:productId", controllers.DeleteProduct) // delete a product by Id

}

func GetIndex(c *fiber.Ctx) error {
	return c.Status(200).JSON(&fiber.Map{
		"message": "Successfully get index page",
	})
}
