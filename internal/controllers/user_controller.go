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
	"strings"
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

func DoTier1(c *fiber.Ctx) error {
	db := c.Locals("db").(*mongo.Database)
	validate := validator.New()

	// Parse multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error parsing form data: " + err.Error(),
		})
	}

	// Extract fields from form data
	doTier := &models.DoTier1{
		WalletAddress: form.Value["wallet_address"][0],
		Bio:           form.Value["bio"][0],
		XLink:         form.Value["x_link"][0],
	}

	// Log parsed struct
	fmt.Printf("Parsed Request Data: %+v\n", doTier)

	// Validate request data
	if err := validate.Struct(doTier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "error parsing form " + err.Error(),
		})
	}

	files := form.File["avatar"]
	var avatarURL string
	userID := doTier.WalletAddress
	fmt.Println("userID is - ", userID)

	for _, fileHead := range files {
		file, err := fileHead.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		fileExtension := filepath.Ext(fileHead.Filename)
		uniqueName := uuid.New().String() + fileExtension

		fmt.Println("avatar bucket is ", os.Getenv("AVATAR_BUCKET"))
		bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to create GridFS bucket: " + err.Error(),
			})
		}

		uploadStream, err := bucket.OpenUploadStream(uniqueName, options.GridFSUpload().SetMetadata(fiber.Map{
			"user_id": userID,
			"ext":     fileExtension,
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
		fmt.Println(avatarURL)
	}

	fmt.Println("doTier.XLink - ", doTier.XLink)
	fmt.Println("doTier.Bio - ", doTier.Bio)

	user := &models.User{
		WalletAddress: userID,
		Avatar:        avatarURL,
		Bio:           doTier.Bio,
		XLink:         doTier.XLink,
	}

	res, err := db.Collection(os.Getenv("USER_COLLECTION")).InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error saving your profile " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "welcome to tier 1",
		"user":    res.InsertedID,
	})
}

func UUserProfile(c *fiber.Ctx) error {
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

	// Fetch user's avatar
	var avatar bson.M
	err = db.Collection("fs.files").FindOne(c.Context(), bson.M{"metadata.user_id": walletAddress}).Decode(&avatar)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			avatarURL = "https://as1.ftcdn.net/v2/jpg/03/46/83/96/1000_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg" // or some default avatar URL
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Error fetching avatar " + err.Error(),
			})
		}
	} else {
		fileID := avatar["_id"].(primitive.ObjectID).Hex()
		baseURL := os.Getenv("BASE_URL")
		avatarURL = fmt.Sprintf("%s/api/v1/user/avatar/%s", baseURL, fileID)
	}

	// Fetch user's products
	var userProducts []models.Product
	cursor, err := db.Collection(os.Getenv("PRODUCT_COLLECTION")).Find(c.Context(), bson.M{"user_id": walletAddress})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error fetching this user's products",
			"details": err.Error(),
		})
	}
	defer cursor.Close(c.Context())

	if !cursor.Next(c.Context()) {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error retrieving user products",
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

func UUpdateUserAvatar(c *fiber.Ctx) error {
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

	if fileHeader.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File size too large, max 5MB allowed",
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
			"message": "error opening file " + err.Error(),
		})
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error reading content, try again " + err.Error(),
		})
	}

	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error reaching bucket, try again " + err.Error(),
		})
	}

	// Check if user already has an avatar
	var existingAvatarMetadata bson.M
	err = db.Collection(os.Getenv("AVATAR_COLLECTION")).FindOne(c.Context(), fiber.Map{"metadata.user_id": user}).Decode(&existingAvatarMetadata)
	fmt.Println(user)
	if err == nil {
		// Delete existing avatar file
		if err := bucket.Delete(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "error deleting existing avatar " + err.Error(),
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
			"message": "error opening upload stream, try again " + err.Error(),
		})
	}

	uploadStream.FileID = user
	defer uploadStream.Close()

	_, err = uploadStream.Write(content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error at upload stream, try again " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Avatar updated successfully",
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

	if fileHeader.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File size too large, max 5MB allowed",
		})
	}

	fileExtension := strings.ToLower(filepath.Ext(fileHeader.Filename))

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
			"message": "error opening file " + err.Error(),
		})
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error reading content, try again " + err.Error(),
		})
	}

	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(os.Getenv("AVATAR_BUCKET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error reaching bucket, try again " + err.Error(),
		})
	}

	// Check if user already has an avatar
	var existingAvatar bson.M
	err = db.Collection("fs.files").FindOne(c.Context(), bson.M{"metadata.user_id": user}).Decode(&existingAvatar)
	if err == nil {
		// Delete existing avatar file
		existingAvatarID := existingAvatar["_id"].(primitive.ObjectID)
		if err := bucket.Delete(existingAvatarID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "error deleting existing avatar " + err.Error(),
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
			"message": "error opening upload stream, try again " + err.Error(),
		})
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "error at upload stream, try again " + err.Error(),
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
