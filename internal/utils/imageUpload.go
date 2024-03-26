package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ImageUploadMiddleware() gin.HandlerFunc {
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
