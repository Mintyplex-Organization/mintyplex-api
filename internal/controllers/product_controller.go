package controllers

import (
	"context"
	"fmt"
	"mintyplex-api/internal/models"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProduct(c *fiber.Ctx) error {
	// Get user ID from params
	user := c.Params("id")

	// Initiate db instance
	db := c.Locals("db").(*mongo.Database)

	// Initiate validator for fields
	validate := validator.New()

	// Parse request to model
	addProd := &models.AddProduct{}
	if err := c.BodyParser(addProd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request data " + err.Error(),
		})
	}

	// Validate request data
	if err := validate.Struct(addProd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Validation error: " + err.Error(),
		})
	}

	fmt.Printf("Parsed Request Data: %+v\n", addProd)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error parsing form data: " + err.Error(),
		})
	}
	files := form.File["image"]

	var imageURL string
	productID := primitive.NewObjectID()

	// Upload image to Cloudinary
	if len(files) > 0 {
		fileHead := files[0]
		file, err := fileHead.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		urlCloudinary := os.Getenv("CLOUDINARY_URL")
		ctx := context.Background()
		cldService, err := cloudinary.NewFromURL(urlCloudinary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to initialize Cloudinary service: " + err.Error(),
			})
		}

		resp, err := cldService.Upload.Upload(ctx, file, uploader.UploadParams{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to upload image to Cloudinary: " + err.Error(),
			})
		}

		imageURL = resp.SecureURL
	}

	product := &models.Product{
		ID:          productID,
		UserId:      user,
		Name:        addProd.Name,
		Price:       addProd.Price,
		Discount:    addProd.Discount,
		Description: addProd.Description,
		Categories:  addProd.Categories,
		Quantity:    addProd.Quantity,
		Tags:        addProd.Tags,
		CoverImage:  imageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	response, err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).InsertOne(c.Context(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal Server Error When Trying To Insert Product Details" + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Product Created successfully",
		"product": response.InsertedID,
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
			"message": "404 - Product Not Found",
			"data":    err.Error(),
		})
	}

	fmt.Println(product.UserId)
	// Fetch user details using the UserId from the product
	var user models.User
	if err := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(ctx, bson.M{"_id": product.UserId}).Decode(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "User Not Found",
			"data":    err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Product and user details fetched successfully",
		"data": fiber.Map{
			"Product": fiber.Map{
				"ID":          product.ID.Hex(),
				"UserId":      product.UserId,
				"CoverImage":  product.CoverImage,
				"Name":        product.Name,
				"Price":       product.Price,
				"Discount":    product.Discount,
				"Description": product.Description,
				"Categories":  product.Categories,
				"Quantity":    product.Quantity,
				"Tags":        product.Tags,
				"CreatedAt":   product.CreatedAt.Format(time.RFC3339),
				"UpdatedAt":   product.UpdatedAt.Format(time.RFC3339),
			},
			"User": fiber.Map{
				"WalletAddress": user.WalletAddress,
				"Avatar":        user.Avatar,
				"Bio":           user.Bio,
				"XLink":         user.XLink,
				"CreatedAt":     user.CreatedAt.Format(time.RFC3339),
				"UpdatedAt":     user.UpdatedAt.Format(time.RFC3339),
			},
		},
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
