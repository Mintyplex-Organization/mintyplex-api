package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddProduct(c *fiber.Ctx) error {
	user := c.Params("id")

	validate := validator.New()
	db := c.Locals("db").(*mongo.Database)

	addProd := &models.AddProduct{}
	if err := c.BodyParser(addProd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request data " + err.Error(),
		})
	}

	if err := validate.Struct(addProd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "validation error: " + err.Error(),
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error parsing form data: " + err.Error(),
		})
	}
	files := form.File["image"]

	// var uploadedFiles []string

	var imageURL string
	productID := primitive.NewObjectID()

	for _, fileHeadr := range files {
		file, err := fileHeadr.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		fileExtension := filepath.Ext(fileHeadr.Filename)
		uniqueFilename := uuid.New().String() + fileExtension
		// uploadedFiles = append(uploadedFiles, uniqueFilename)

		bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("COVER_BUCKET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to create GridFS bucket: " + err.Error(),
			})
		}

		uploadStream, err := bucket.OpenUploadStream(uniqueFilename, options.GridFSUpload().SetMetadata(fiber.Map{
			"product_id": productID,
			"ext":        fileExtension,
		}))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to open GridFS upload stream: " + err.Error(),
			})
		}
		defer uploadStream.Close()

		if _, err := io.Copy(uploadStream, file); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to copy file data to GridFS upload stream: " + err.Error(),
			})
		}

		imageURL = fmt.Sprintf("%s/api/v1/product/cover/%s", os.Getenv("BASE_URL"), productID.Hex())
	}

	fmt.Println(imageURL)

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
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	fmt.Println(product.CoverImage)
	fmt.Println(product.Quantity)
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
		"task":    response.InsertedID,
	})
}

func GetProductCover(c *fiber.Ctx) error {
	productID := c.Params("id")
	fmt.Println("Product ID:", productID)

	var coverMetadata bson.M

	db := c.Locals("db").(*mongo.Database)

	// Convert the product ID string to an ObjectId
	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		fmt.Println("Invalid product ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid product ID",
			"data":    err.Error(),
		})
	}

	// Query the document using the product ID
	if err := db.Collection(os.Getenv("COVER_COLLECTION")).FindOne(c.Context(), bson.M{"metadata.product_id": objectID}).Decode(&coverMetadata); err != nil {
		fmt.Println("Error fetching cover metadata:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Cover metadata retrieval error",
			"data":    err.Error(),
		})
	}

	fmt.Println("Cover Metadata:", coverMetadata)

	var buffer bytes.Buffer
	bucket, _ := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("COVER_BUCKET")))
	_, err = bucket.DownloadToStreamByName(coverMetadata["filename"].(string), &buffer)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error downloading image",
			"data":    err.Error(),
		})
	}

	// Set response headers
	utils.SetAvatarHeaders(c, buffer, coverMetadata["metadata"].(bson.M)["ext"].(string))

	// Return the image data as the response body
	return c.Send(buffer.Bytes())
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

	baseURL := os.Getenv("BASE_URL")
	product.CoverImage = fmt.Sprintf("%s%s", baseURL, product.CoverImage)

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
