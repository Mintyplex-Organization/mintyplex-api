package controllers

import (
	"mintyplex-api/internal/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// func ReserveUsername(c *fiber.Ctx) error {
// 	reserveuname := new(models.ReserveUsername)
// 	// c.BodyParser(&reserveuname)
// 	if err := c.BodyParser(&reserveuname); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Invalid operation",
// 			"data":    err,
// 		})
// 	}

// 	matched, err := regexp.MatchString(`^[a-zA-Z]+$`, reserveuname.Username)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Username will take only alphabets",
// 			"data":    err,
// 		})
// 	}
// 	if !matched {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Username can only contain alphabets",
// 		})
// 	}

// 	uname := &models.Reserved{}

// 	uname.Username = reserveuname.Username
// 	uname.Email = reserveuname.Email
// 	uname.CreatedAt = time.Now()
// 	uname.UpdatedAt = time.Now()

// 	db := c.Locals("db").(*mongo.Database)

// 	if err := db.Collection(os.Getenv("USERNAME_COLLECTION")).FindOne(c.Context(), fiber.Map{"username": uname.Username}).Decode(&uname); err == nil {
// 		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "username/email is taken",
// 		})
// 	}

// 	if err := db.Collection(os.Getenv("USERNAME_COLLECTION")).FindOne(c.Context(), fiber.Map{"email": uname.Email}).Decode(&uname); err == nil {
// 		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "email/username is taken",
// 		})
// 	}

// 	if _, err := db.Collection(os.Getenv("USERNAME_COLLECTION")).InsertOne(c.Context(), uname); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "invalid or repeating credentials" + err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"error":   false,
// 		"message": "successfully reserved username " + uname.Username,
// 	})
// }

func ReserveUsername(c *fiber.Ctx) error {
	// Parse the request body into the ReserveUsername model
	reserveuname := new(models.ReserveUsername)
	if err := c.BodyParser(&reserveuname); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid operation",
			"data":    err.Error(),
		})
	}

	// Create a new Reserved instance and populate it with the request data
	uname := &models.Reserved{
		Username:  reserveuname.Username,
		Email:     reserveuname.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Get the database instance from the context
	db := c.Locals("db").(*mongo.Database)

	// Check if the username is already taken
	var existingUsername models.Reserved
	if err := db.Collection(os.Getenv("USERNAME_COLLECTION")).FindOne(c.Context(), bson.M{"username": uname.Username}).Decode(&existingUsername); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Username is already taken",
		})
	}

	// Check if the email is already taken
	var existingEmail models.Reserved
	if err := db.Collection(os.Getenv("USERNAME_COLLECTION")).FindOne(c.Context(), bson.M{"email": uname.Email}).Decode(&existingEmail); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Email is already taken",
		})
	}

	// Insert the new username and email into the database
	if _, err := db.Collection(os.Getenv("USERNAME_COLLECTION")).InsertOne(c.Context(), uname); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid or repeating credentials: " + err.Error(),
		})
	}

	// Respond with success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Successfully reserved username " + uname.Username,
	})
}
