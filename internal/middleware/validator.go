package middleware

import (
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/utils"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func ValidateJwt() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		user := &models.User{}

		claims, err := utils.ExtractTokenMetadata(ctx)
		if err != nil {
			// Return status 500 and JWT parse error.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		if claims.Expires < time.Now().Unix() {
			// Return status 401 and JWT expired error.
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Token expired",
			})
		}

		db := ctx.Locals("db").(*mongo.Database)

		if err := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(ctx.Context(), fiber.Map{"_id": claims.UserID}).Decode(&user); err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Invalid token",
			})
		}

		bearToken := strings.Split(ctx.Get("Authorization"), " ")[1]
		tokens := user.Tokens
		tokenExists := false

		for _, token := range tokens {
			if token == bearToken {
				tokenExists = true
				break
			}
		}

		if !tokenExists {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Token does not exist",
			})
		}

		ctx.Locals("user", user)
		return ctx.Next()

	}
}
