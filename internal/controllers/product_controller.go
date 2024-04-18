package controllers

import (
	"context"
	"fmt"
	"io"
	"mintyplex-api/internal/models"
	"net/http"
	"os"
	"strings"
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

// func modItemUpload(c *fiber.Ctx) error {
// 	type RequestBody struct {
// 		WalletAddress string `json:"wallet_address"`
// 		Type string `json:"type"`
// 		File interface{} `json:"file"`
// 	}

// 	var reqBody RequestBody
// 	err := c.BodyParser(&reqBody)
// 	if err != nil{
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": true,
// 			"message": "request is invalid",
// 		})
// 	}

// 	if reqBody.Type == "ebook" && !(strings.HasSuffix(reqBody.File.(string), ".pdf") || strings.HasSuffix(reqBody.File.(string), ".epub")){
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			  "error": true,
//       "message": "Invalid ebook format. Only .pdf and .epub extensions allowed",
// 		})
// 	}

// 	db := c.Locals("db").(*mongo.Database)

// 	var user models.User
// 	err = db.Collection(os.Getenv("USER_COLLECTION")).FindOne(c.Context(), bson.M{"_id": reqBody.WalletAddress}).Decode(user)
// 	if err != nil{
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": true,
// 			"message": "incomplete request",
// 		})
// 	}

// 	baseUrl := "http://localhost:9980/api/worker/objects/"
// 	itemPath := fmt.Sprintf("%s/%s/%s", reqBody.WalletAddress, reqBody.Type, "[nameofitem]") // Replace with actual file name handling
// 	url := baseUrl+itemPath
// 	fmt.Println(url)


// }

func ItemUpload(c *fiber.Ctx) error {
	// Authorization token or credentials
	authToken := "Basic UGFzc3dvcmQ6MjQ3YWRtaW5pc3RyYXRpb24="

	url := "http://localhost:9980/api/worker/objects/user_products/sion419619/video_the-essense-of-work"
	method := "PUT"

	payload := strings.NewReader("<file contents here>")
	fmt.Println(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
 
	// Add authorization header
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	fmt.Println(string(body))
	fmt.Println(body)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Item uploaded successfully",
	})
}

// func UpdateProduct(c *fiber.Ctx)error{
// 	ctx, cancel := context.WithTimeout(context.TODO(), 50*time.Second)
// 	defer cancel()

// 	var updateData bson.M
// 	if err := c.BodyParser(&updateData); err != nil{
// 		return c.Status(400).JSON(fiber.Map{
// 			"status": "error",
// 			"message": "please refresh page and try again",
// 			"data": nil,
// 		})
// 	}

// }
