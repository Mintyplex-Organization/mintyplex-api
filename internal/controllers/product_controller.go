package controllers

import (
	"context"
	"fmt"
	"mintyplex-api/internal/models"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProduct(c *fiber.Ctx) error {
	// user := c.Locals("user").(*models.User)
	user := c.Params("id")

	// userID := c.Params("id") // Assuming "id" parameter is the user's ID

	fmt.Println(user)

	validate := validator.New()
	db := c.Locals("db").(*mongo.Database)

	var usr models.User
	err := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(c.Context(), bson.M{"_id": user}).Decode(&usr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request",
		})
	}

	addProduct := new(models.AddProduct)
	if err := c.BodyParser(&addProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if err := validate.Struct(addProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	product := new(models.Product)

	timestamp := time.Now().Unix()

	product.UserId = user
	product.Name = addProduct.Name
	product.Price = addProduct.Price
	product.Discount = addProduct.Discount
	product.Description = addProduct.Description
	product.Categories = addProduct.Categories
	product.Quantity = addProduct.Quantity
	product.Tags = addProduct.Tags
	product.CreatedAt = timestamp
	product.UpdatedAt = timestamp

	response, err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).InsertOne(c.Context(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Product Created successfully",
		"task":    response.InsertedID,
	})

}

func AllProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 50*time.Second)
	defer cancel()

	db := c.Locals("db").(*mongo.Database)
	cursor, err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).Find(ctx, bson.M{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "error getting products",
			"data":    err,
		})
	}

	var products []bson.M
	if err := cursor.All(ctx, &products); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error parsing products",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Products fetched successfully",
		"data":    products,
	})

}
