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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func OneProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 50*time.Second)
	defer cancel()

	id := c.Params("id")

	productId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product id",
			"data":    err,
		})
	}

	var product models.Product
	db := c.Locals("db").(*mongo.Database)
	if err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).FindOne(ctx, bson.M{"_id": productId}).Decode(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "link is broken",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Product fetched successfully",
		"data":    product,
	})

}

func UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	userID := c.Params("uid")

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid product ID",
		})
	}

	filter := bson.M{"_id": productObjectID, "user_id": userID}

	db := c.Locals("db").(*mongo.Database)

	var existingProduct models.Product
	if err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).FindOne(c.Context(), filter).Decode(&existingProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "product does not exist",
		})
	}
	var UpdateProduct models.Product
	if err := c.BodyParser(&UpdateProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid update data for model binding, please refer to customer care",
			"data":    err,
		})
	}

	if existingProduct.UserId != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid reach, try again",
		})
	}

	update := bson.M{"$set": UpdateProduct}
	productCollection := db.Collection(os.Getenv("PRODUCT_COLLECTION"))
	if _, err := productCollection.UpdateByID(c.Context(), productObjectID, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error updating product, we can try again",
			"data":    err,
		})
	}

	if err := productCollection.FindOne(c.Context(), filter).Decode(&existingProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "error reaching product",
			"data":    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Product updated successfully",
		"data":    existingProduct,
	})
}

func UUpdateProduct(c *fiber.Ctx) error {
	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "request is invalid",
			"data":    err,
		})
	}

	id := c.Params("id")

	productId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "link is broken, try again",
			"data":    err,
		})
	}

	filter := bson.M{"_id": productId}

	db := c.Locals("db").(*mongo.Database)
	product, err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).FindOneAndUpdate(c.Context(), filter, bson.M{"$set": updateData}).DecodeBytes()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error updating product",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Product updated successfully",
		"data":    product.String(),
	})
}

// func DeleteProduct(c *fiber.Ctx) error {

// }
