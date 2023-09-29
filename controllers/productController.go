package controllers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/BagusNurhudaUI/dagangan-bagus-nurhuda/models"

	"github.com/gofiber/fiber/v2"
)

// function for getting all the products
func (db *InDB) GetProduct(c *fiber.Ctx) error {

	Product := []models.Product{}

	err := db.DB.Debug().Find(&Product).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	} else {
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"data": Product,
		})
	}

}

// functions for getting product by ID
func (db *InDB) GetProductById(c *fiber.Ctx) error {

	params, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	Product := models.Product{}

	err = db.DB.Debug().Where("id = ?", params).First(&Product).Error
	log.Println(Product.ToString())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	} else {
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"data": Product,
		})
	}

}

// functions for post a products
func (db *InDB) PostProduct(c *fiber.Ctx) error {

	Product := new(models.Product)

	// parse the req.body to product models
	if err := c.BodyParser(Product); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	//create a new product
	err := db.DB.Debug().Create(&Product).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	} else {
		return c.Status(http.StatusCreated).JSON(&fiber.Map{
			"message": "Product successfully created",
			"data":    Product,
		})
	}

}

// functions for update the product by id, fields (title, caption, photo_url and price) will be updated
func (db *InDB) UpdateProduct(c *fiber.Ctx) error {
	params := c.Params("productId")
	Product := new(models.Product)       // models to get a fields product from database
	UpdateProduct := new(models.Product) // models to update a field product to database

	// parse req.body to product models
	if err := c.BodyParser(Product); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	//find id product
	err := db.DB.Debug().Where("id = ?", params).First(&UpdateProduct)
	if err.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": fmt.Sprintf("Update not successful, id %s is not found", params),
		})

	} else {
		if Product.Title != "" {
			UpdateProduct.Title = Product.Title
		}
		if Product.Caption != "" {
			UpdateProduct.Caption = Product.Caption
		}
		if Product.Photo_url != "" {
			UpdateProduct.Photo_url = Product.Photo_url
		}
		if Product.Price != 0 {
			UpdateProduct.Price = Product.Price
		}

		// Update the product from request body
		err := db.DB.Debug().Save(&UpdateProduct).Error
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": err.Error(),
			})
		} else {
			return c.Status(http.StatusOK).JSON(&fiber.Map{
				"message": "Updated successfully",
				"data":    UpdateProduct,
			})
		}

	}

}

// function to delete a product by id
func (db *InDB) DeleteProduct(c *fiber.Ctx) error {
	params := c.Params("productId")
	Product := models.Product{}

	// Delete the product by 1
	err := db.DB.Debug().Delete(&Product, params)
	if err.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error.Error(),
		})
	} else if err.RowsAffected < 1 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": fmt.Sprintf("Delete not successful, id %s is not found", params),
		})
	} else {
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Your product has been successfully deleted",
		})
	}
}

// function to get pagination information from products
func (db *InDB) Paginate(c *fiber.Ctx) error {

	Product := []models.Product{}
	rawQuery := "SELECT * FROM products"
	var (
		rawSearch string
		rawPrice  string
	)
	search := c.Query("search")
	maxPrice, _ := strconv.Atoi(c.Query("maxprice"))
	minPrice, _ := strconv.Atoi(c.Query("minprice"))
	sortbyprice := c.Query("sortbyprice")

	// assign query parameters SEARCH if exists to raw sql
	if search != "" {
		rawSearch = fmt.Sprintf("(Title LIKE '%%%s%%' OR Caption LIKE '%%%s%%')", search, search)
	}

	// assign query parameters PRICE if exists to raw sql
	if maxPrice != 0 && minPrice != 0 {
		rawPrice = fmt.Sprintf("(price <= %d AND price >= %d)", maxPrice, minPrice)
	} else {
		if maxPrice != 0 {
			rawPrice = fmt.Sprintf("%s price <= %d", rawPrice, maxPrice)
		}

		if minPrice != 0 {
			rawPrice = fmt.Sprintf("%s price >= %d", rawPrice, minPrice)
		}
	}

	// Merge the rawPrice and rawSearch if exists to rawQuery
	if rawPrice != "" && rawSearch != "" {
		rawQuery = fmt.Sprintf("%s WHERE %s AND %s", rawQuery, rawSearch, rawPrice)
	} else {
		if rawSearch != "" {
			rawQuery = fmt.Sprintf("%s WHERE %s ", rawQuery, rawSearch)
		}

		if rawPrice != "" {
			rawQuery = fmt.Sprintf("%s WHERE %s ", rawQuery, rawPrice)
		}
	}

	// Make the sort query string
	if sortbyprice != "" && (sortbyprice == "true" || sortbyprice == "false") {
		if sortbyprice == "true" {
			rawQuery = fmt.Sprintf("%s ORDER BY price ASC", rawQuery)
		}
		if sortbyprice == "false" {
			rawQuery = fmt.Sprintf("%s ORDER BY price DESC", rawQuery)
		}
	} else {
		rawQuery = fmt.Sprintf("%s ORDER BY id ASC", rawQuery)
	}

	

	// Get the page and the limit for the paginations.
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "5"))

	temp := []models.Product{}
	// Get the total number of rows in the database Product
	db.DB.Raw(rawQuery).Scan(&temp)
	totalRows := len(temp)
	totalPage := math.Ceil(float64(float64(totalRows)/float64(limit)) + 0.00000000001)

	if page > int(totalPage) {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Page is out of bounds, not found",
		})
	}
	rawQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", rawQuery, int64(limit), limit*(page-1))
	rawQuery = fmt.Sprintf("%s ;", rawQuery)

	// Get the Product By Final Queries
	err := db.DB.Debug().Raw(rawQuery).Scan(&Product).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"search":  search,
			"message": err.Error(),
		})
	} else {
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"data":      Product,
			"totalRows": totalRows,
			"totalPage": totalPage,
			"page":      page,
			"limit":     limit,
		})
	}
}
