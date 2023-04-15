package router

import (
	"github.com/BagusNurhudaUI/dagangan-bagus-nurhuda/controllers"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func StartApp(app *fiber.App, db *gorm.DB) {
	InDB := controllers.New(db)

	app.Get("/", GetIndex)

	products := app.Group("/products")
	products.Get("/", InDB.Paginate)                   // get all products with pagination
	products.Get("/all", InDB.GetProduct)              // get all products without pagination
	products.Get("/:productId", InDB.GetProductById)   // get all products by productId
	products.Post("/", InDB.PostProduct)               // create a new product
	products.Put("/:productId", InDB.UpdateProduct)    // update a product with existing data
	products.Delete("/:productId", InDB.DeleteProduct) // delete a product by Id

}

func GetIndex(c *fiber.Ctx) error {
	return c.Status(200).JSON(&fiber.Map{
		"message": "Successfully get index page",
	})
}
