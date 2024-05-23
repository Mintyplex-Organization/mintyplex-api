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
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/internal/uuid"
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



func DoTier91(c *fiber.Ctx) error {
	db := c.Locals("db").(*mongo.Database)

	validate := validator.New()

	dotier1 := &models.DoTier1{}
	if err := c.BodyParser(dotier1); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "the from cannot be processed now " + err.Error(),
		})
	}

	// validate request data
	if err := validate.Struct(dotier1); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "validation failed for this form" + err.Error(),
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error parsing form data: " + err.Error(),
		})
	}
	avatar := form.File["avatar"]

	var avatarURL string
	userID := dotier1.WalletAddress

	for _, avatarHeader := range avatar {
		file, err := avatarHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		fileExtension := filepath.Ext(avatarHeader.Filename)
		uniqueFileName := uuid.New().String() + fileExtension

		bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to create GridFS bucket: " + err.Error(),
			})
		}

		uploadStream, err := bucket.OpenUploadStream(uniqueFileName, options.GridFSUpload().SetMetadata(fiber.Map{
			"avatar_id": userID,
			"ext":       fileExtension,
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

		avatarURL = fmt.Sprintf("%s/api/v1/user/avatar/%s", os.Getenv("BASE_URL"), userID)
	}

	user := &models.User{
		WalletAddress: dotier1.WalletAddress,
		XLink:         dotier1.XLink,
		Bio:           dotier1.Bio,
		Avatar:        avatarURL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	res, err := db.Collection(os.Getenv("USER_COLLECTION")).InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Internal Server Error When Trying To Insert User Details" + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Product Created successfully",
		"user":    res.InsertedID,
	})
}

func DooTier1(c *fiber.Ctx) error {
	// Initiate DB instance
	db := c.Locals("db").(*mongo.Database)

	// Parse request to model
	dotier1 := new(models.DoTier1)
	if err := c.BodyParser(dotier1); err != nil {
		fmt.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
	}

	fmt.Println("Parsed request:", dotier1)

	// Ensure WalletAddress is correctly parsed
	if dotier1.WalletAddress == "" {
		fmt.Println("Wallet address is missing in the request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Wallet address is required",
		})
	}

	// Check if user already exists
	var existingUser models.User
	if err := db.Collection(os.Getenv("USER_COLLECTION")).FindOne(c.Context(), bson.M{"wallet_address": dotier1.WalletAddress}).Decode(&existingUser); err == nil {
		fmt.Println("User already exists:", existingUser)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "This wallet cannot have another account",
		})
	}

	// Create new user ID
	userID := dotier1.WalletAddress
	fmt.Println("User ID is:", userID)

	// Handle avatar upload
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Error parsing multipart form:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error parsing form data: " + err.Error(),
		})
	}

	fmt.Println("Parsed multipart form:", form)

	files := form.File["avatar"]
	if len(files) == 0 {
		fmt.Println("No avatar file provided")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No avatar file provided",
		})
	}

	var avatarURL string

	for _, fileHeader := range files {
		fmt.Println("Processing file:", fileHeader.Filename)

		if fileHeader.Size > 5*1024*1024 { // Increase file size limit to 5MB
			fmt.Println("File size too large:", fileHeader.Size)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "File size too large, max 5MB allowed",
			})
		}

		fileExtension := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".png" {
			fmt.Println("Invalid file type:", fileExtension)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid file type",
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error: " + err.Error(),
			})
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading file content:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error: " + err.Error(),
			})
		}

		bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
		if err != nil {
			fmt.Println("Error creating GridFS bucket:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error: " + err.Error(),
			})
		}

		uploadStream, err := bucket.OpenUploadStream(userID, options.GridFSUpload().SetMetadata(bson.M{
			"user_id": userID,
			"ext":     fileExtension,
		}))
		if err != nil {
			fmt.Println("Error opening upload stream:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error: " + err.Error(),
			})
		}
		defer uploadStream.Close()

		if _, err := uploadStream.Write(content); err != nil {
			fmt.Println("Error writing to upload stream:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error: " + err.Error(),
			})
		}

		avatarURL = fmt.Sprintf("%s/api/v1/user/avatar/%s", os.Getenv("BASE_URL"), userID)
	}

	fmt.Println("Avatar URL is:", avatarURL)

	// Create new user
	user := &models.User{
		WalletAddress: userID,
		XLink:         dotier1.XLink,
		Bio:           dotier1.Bio,
		Avatar:        avatarURL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	fmt.Println("New user data:", user)

	// Insert the new user into the database
	if _, err := db.Collection(os.Getenv("USER_COLLECTION")).InsertOne(c.Context(), user); err != nil {
		fmt.Println("Error inserting new user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create user: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Successfully upgraded to Tier 1",
	})
}

func UserProfile(c *fiber.Ctx) error {
	walletAddress := c.Params("id")

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

	var avatarURL string

	avatarID := foundUser.WalletAddress
	if avatarID != "" {
		avatarURL = "https://mintyplex-api.onrender.com/api/v1/user/avatar/" + avatarID
	}

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
	prodColl := db.Collection(os.Getenv("PRODUCT_COLLECTION"))

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "error getting users",
			"data":    err,
		})
	}

	defer cursor.Close(ctx)

	var usersProducts []map[string]interface{}

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "error decoding user",
				"data":    err.Error(),
			})
		}

		prodCur, err := prodColl.Find(ctx, bson.M{"user_id": user.WalletAddress})
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "error getting user products for " + user.WalletAddress,
				"data":    err.Error(),
			})
		}
		defer prodCur.Close(ctx)

		var userProd []models.Product
		if err := prodCur.All(ctx, &userProd); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "error fetching products for user " + user.WalletAddress,
				"data":    err.Error(),
			})
		}

		userData := map[string]interface{}{
			"userProfile": user,
			"products":    userProd,
		}
		usersProducts = append(usersProducts, userData)

	}

	// var users []models.User
	// if err := cursor.All(ctx, &users); err != nil {
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "error fetching users",
	// 		"data":    err,
	// 	})
	// }

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "users retrieved",
		"data":    usersProducts,
	})
}
