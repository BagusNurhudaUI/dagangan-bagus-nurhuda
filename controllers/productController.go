package controllers

import (
	"dagangan/config"
	"dagangan/models"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProduct(c *fiber.Ctx) error {
	db := config.GetDB()
	Product := []models.Product{}

	err := db.Debug().Find(&Product).Error
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

func GetProductById(c *fiber.Ctx) error {
	db := config.GetDB()
	params := c.Params("productId")
	Product := models.Product{}
	_ = db

	err := db.Debug().Where("id = ?", params).First(&Product).Error
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

func PostProduct(c *fiber.Ctx) error {
	db := config.GetDB()
	// contentType := c.Request.Header.Get("Content-Type")
	Product := new(models.Product)
	_ = db
	if err := c.BodyParser(Product); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := db.Debug().Create(&Product).Error
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

func UpdateProduct(c *fiber.Ctx) error {
	db := config.GetDB()
	params := c.Params("productId")
	_, _ = db, params
	Product := new(models.Product)
	UpdateProduct := new(models.Product)

	if err := c.BodyParser(Product); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := db.Debug().Where("id = ?", params).First(&UpdateProduct)
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
		err := db.Debug().Save(&UpdateProduct).Error
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

func DeleteProduct(c *fiber.Ctx) error {
	db := config.GetDB()
	params := c.Params("productId")
	Product := []models.Product{}
	err := db.Debug().Delete(&Product, params)
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

func Paginate(c *fiber.Ctx) error {
	db := config.GetDB()
	Product := []models.Product{}
	rawQuery := "SELECT * FROM products"
	var (
		rawSearch string
		// totalRows int64
		rawPrice string
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
	if sortbyprice != "" {
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
	db.Raw(rawQuery).Scan(&temp)
	totalRows := len(temp)

	totalPage := math.Ceil(float64(float64(totalRows)/float64(limit)) + 0.00000000001)

	if page > int(totalPage) {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Page is out of bounds, not found",
		})
	}
	rawQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", rawQuery, int64(limit), limit*(page-1))

	// Get the Product By Final Queries
	err := db.Debug().Raw(rawQuery).Scan(&Product).Error
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
