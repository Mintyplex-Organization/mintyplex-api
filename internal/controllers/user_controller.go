package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/utils"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DoTier1(c *fiber.Ctx) error {

	dotier1 := new(models.DoTier1)
	c.BodyParser(&dotier1)

	user := &models.User{}

	user.WalletAddress = dotier1.WalletAddress
	user.XLink = dotier1.XLink
	user.Bio = dotier1.Bio
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	db := c.Locals("db").(*mongo.Database)

	if err := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(c.Context(), fiber.Map{"wallet_address": user.WalletAddress}).Decode(&user); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "this wallet cannot have another account",
		})
	}

	if _, err := db.Collection(os.Getenv("USER_COLLECTION")).InsertOne(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "invalid and repeating credentials" + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "successfully upgraded to Tier 1",
	})
}

func UserProfile(c *fiber.Ctx) error {
	walletAddress := c.Params("id")
	// baseURL := "https://mintyplex-api.onrender.com"

	db := c.Locals("db").(*mongo.Database)

	result := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(c.Context(), fiber.Map{"_id": walletAddress})
	var err error
	if result.Err() != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "User not found, give it a second",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "User not found rounds, give it a second",
			})
		}

	}

	var foundUser models.User
	if err := result.Decode(&foundUser); err != nil {
		return err
	}

	// avatarURL := baseURL + "/api/v1/user/avatar/" + foundUser.WalletAddress
	var avatarURL string
	avatarID := foundUser.WalletAddress
	if avatarID != "" {
		avatarURL = "https://mintyplex-api.onrender.com/api/v1/user/avatar/" + avatarID
	}

	fmt.Println(avatarID)
	fmt.Println(avatarURL)

	var userProducts []models.Product
	cursor, err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).Find(c.Context(), bson.M{"user_id": walletAddress})
	fmt.Println(cursor)
	fmt.Println(walletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error fetching this user's products",
			"details": err.Error(),
		})

	}
	defer cursor.Close(c.Context())

	if cursor.Next(c.Context()) != true {
		return c.Status(200).JSON(fiber.Map{
			"error":   false,
			"message": "User Profile",
			"user": models.UserProfileResponse{
				WalletAddress: foundUser.WalletAddress,
				XLink:         foundUser.XLink,
				Bio:           foundUser.Bio,
				Avatar:        avatarURL,
				Products:      []models.Product{},
				CreatedAt:     foundUser.CreatedAt,
				UpdatedAt:     foundUser.UpdatedAt,
			},
		})

	}
	fmt.Println("x_link ", foundUser.XLink)

	if err := cursor.All(c.Context(), &userProducts); err != nil {
		fmt.Println(&userProducts)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error retreiving user products",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "User Profile",
		"user": models.UserProfileResponse{
			WalletAddress: foundUser.WalletAddress,
			XLink:         foundUser.XLink,
			Bio:           foundUser.Bio,
			Avatar:        avatarURL,
			Products:      userProducts,
			CreatedAt:     foundUser.CreatedAt,
			UpdatedAt:     foundUser.UpdatedAt,
		},
	})
}

func UpdateUserProfile(c *fiber.Ctx) error {
	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid data, try again",
			"data":    err,
		})
	}

	id := c.Params("id")
	db := c.Locals("db").(*mongo.Database)

	filter := bson.M{"_id": id}

	profile, err := db.Collection(os.Getenv("USER_COLLECTION")).FindOneAndUpdate(c.Context(), filter, bson.M{"$set": updateData}).DecodeBytes()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error updating profile",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "updated!",
		"data":    profile.String(),
	})
}

func UploadUserAvatar(c *fiber.Ctx) error {
	userID := c.Params("id")

	db := c.Locals("db").(*mongo.Database)

	// Fetch user using wallet address
	var user models.User
	err := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "User not found " + err.Error(),
		})
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "file header is not type of avatar " + err.Error(),
		})
	}

	if fileHeader.Size > 1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File size too large, max 1MB allowed",
		})
	}

	fileExtension := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):])

	if fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file type",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	var avatarMetadata bson.M

	if err := db.Collection(os.Getenv("AVATAR_COLLECTION")).FindOne(c.Context(), fiber.Map{"metadata.user_id": userID}).Decode(&avatarMetadata); err == nil {
		// Delete existing avatar file
		if err := bucket.Delete(userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error",
			})
		}
	}

	uploadStream, err := bucket.OpenUploadStream(fileHeader.Filename, options.GridFSUpload().SetMetadata(fiber.Map{
		"user_id": userID,
		"ext":     fileExtension,
	}))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	uploadStream.FileID = userID
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	log.Printf("Write file to DB was successful. File size: %d KB\n", fileSize/1024)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Avatar uploaded successfully",
	})
}

func UpdateUserAvatar(c *fiber.Ctx) error {
	user := c.Params("id")

	db := c.Locals("db").(*mongo.Database)

	// Fetch user using wallet address
	var usr models.User
	err := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(c.Context(), bson.M{"_id": user}).Decode(&usr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "User not found " + err.Error(),
		})
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "file header is not type of avatar " + err.Error(),
		})
	}

	if fileHeader.Size > 1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File size too large, max 1MB allowed",
		})
	}

	fileExtension := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):])

	if fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file type",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	// Check if user already has an avatar
	var existingAvatarMetadata bson.M
	err = db.Collection(os.Getenv("AVATAR_COLLECTION")).FindOne(c.Context(), fiber.Map{"metadata.user_id": user}).Decode(&existingAvatarMetadata)
	if err == nil {
		// Delete existing avatar file
		if err := bucket.Delete(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error",
			})
		}
	}

	uploadStream, err := bucket.OpenUploadStream(fileHeader.Filename, options.GridFSUpload().SetMetadata(fiber.Map{
		"user_id": user,
		"ext":     fileExtension,
	}))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	uploadStream.FileID = user
	defer uploadStream.Close()

	_, err = uploadStream.Write(content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Avatar updated successfully",
	})
}

func GetAvatarById(c *fiber.Ctx) error {
	userID := c.Params("id")

	var avatarMetadata bson.M

	db := c.Locals("db").(*mongo.Database)

	if err := db.Collection(os.Getenv("AVATAR_COLLECTION")).FindOne(c.Context(), fiber.Map{"metadata.user_id": userID}).Decode(&avatarMetadata); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Avatar not found",
		})
	}

	var buffer bytes.Buffer
	bucket, _ := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
	bucket.DownloadToStream(userID, &buffer)

	utils.SetAvatarHeaders(c, buffer, avatarMetadata["metadata"].(bson.M)["ext"].(string))

	return c.Send(buffer.Bytes())
}

func DeleteUserAvatar(c *fiber.Ctx) error {
	userID := c.Params("id")
	db := c.Locals("db").(*mongo.Database)

	var avatarMetadata bson.M

	if err := db.Collection(os.Getenv("AVATAR_COLLECTION")).FindOne(c.Context(), fiber.Map{"metadata.user_id": userID}).Decode(&avatarMetadata); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Avatar not found",
		})
	}

	bucket, _ := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))

	if err := bucket.Delete(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Avatar deleted successfully",
	})
}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 35*time.Second)
	defer cancel()
	db := c.Locals("db").(*mongo.Database)
	collection := db.Collection(os.Getenv("USER_COLLECTION"))
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "error getting users",
			"data":    err,
		})
	}

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "error fetching users",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "users retrieved",
		"data":    users,
	})
}
