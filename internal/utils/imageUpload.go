package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

func MImageUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}
		defer file.Close()

		c.Set("filepath", header.Filename)
		c.Set("file", file)

		c.Next()
	}
}

func ImageUploadMiddleware(c *fiber.Ctx) error {
	// Get the file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	// defer file.Close()

	// Set the filepath and file in the context locals
	c.Locals("filepath", file.Filename)
	c.Locals("file", file)

	return c.Next()
}
