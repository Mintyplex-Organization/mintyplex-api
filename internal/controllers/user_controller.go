package controllers

import (
	"context"
	"fmt"
	"mintyplex-api/internal/models"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DoTier1(c *fiber.Ctx) error {
	db := c.Locals("db").(*mongo.Database)
	validate := validator.New()
	urlCloudinary := os.Getenv("CLOUDINARY_URL")

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "error parsing form data: " + err.Error(),
		})
	}

	dotier1 := &models.DoTier1{
		WalletAddress: form.Value["wallet_address"][0],
		Bio:           form.Value["bio"][0],
		XLink:         form.Value["x_link"][0],
	}

	fmt.Printf("Parsed Request Data: %+v\n", dotier1)

	// Validate request data
	if err := validate.Struct(dotier1); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "error parsing form " + err.Error(),
		})
	}

	files := form.File["avatar"]
	userID := dotier1.WalletAddress
	var avatarURL string

	if len(files) > 0 {
		for _, fileHead := range files {
			file, err := fileHead.Open()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   true,
					"message": "Failed to open image file: " + err.Error(),
				})
			}
			defer file.Close()

			ctx := context.Background()
			cldService, err := cloudinary.NewFromURL(urlCloudinary)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   true,
					"message": "Failed to create Cloudinary service: " + err.Error(),
				})
			}

			resp, err := cldService.Upload.Upload(ctx, file, uploader.UploadParams{})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   true,
					"message": "Failed to upload image to Cloudinary: " + err.Error(),
				})
			}

			avatarURL = resp.SecureURL
			fmt.Println(avatarURL)
		}
	}

	if avatarURL == "" {
		avatarURL = fmt.Sprintf("%s/api/v1/user/avatar/default.png", os.Getenv("BASE_URL"))
	}

	fmt.Println("doTier.XLink - ", dotier1.XLink)
	fmt.Println("doTier.Bio - ", dotier1.Bio)

	user := &models.User{
		WalletAddress: userID,
		Avatar:        avatarURL,
		Bio:           dotier1.Bio,
		XLink:         dotier1.XLink,
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
		"message": "you're welcome to tier 1, have a look around",
		"user":    res.InsertedID,
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

	// var avatarURL string

	// avatarID := foundUser.WalletAddress
	// if avatarID != "" {
	// 	avatarURL = "https://mintyplex-api.onrender.com/api/v1/user/avatar/" + avatarID
	// }

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
				Avatar:        foundUser.Avatar,
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
			Avatar:        foundUser.Avatar,
			Products:      userProducts,
			CreatedAt:     foundUser.CreatedAt,
			UpdatedAt:     foundUser.UpdatedAt,
		},
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

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "users retrieved",
		"data":    usersProducts,
	})
}

func UpdateUserProfile(c *fiber.Ctx) error {
	// Get the user ID from the URL parameters
	id := c.Params("id")
	db := c.Locals("db").(*mongo.Database)

	// Parse form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error parsing form data: " + err.Error(),
		})
	}

	// Extract fields from form data
	updateData := bson.M{}
	if bio := form.Value["bio"]; len(bio) > 0 {
		updateData["bio"] = bio[0]
	}
	if xLink := form.Value["x_link"]; len(xLink) > 0 {
		updateData["x_link"] = xLink[0]
	}

	files := form.File["avatar"]
	if len(files) > 0 {
		// Assuming only one file is uploaded for the avatar
		fileHead := files[0]
		file, err := fileHead.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		// Upload the file to Cloudinary
		urlCloudinary := os.Getenv("CLOUDINARY_URL")
		ctx := context.Background()
		cldService, err := cloudinary.NewFromURL(urlCloudinary)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to initialize Cloudinary service: " + err.Error(),
			})
		}

		resp, err := cldService.Upload.Upload(ctx, file, uploader.UploadParams{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to upload image to Cloudinary: " + err.Error(),
			})
		}

		// Get the secure URL from the Cloudinary response
		avatarURL := resp.SecureURL
		updateData["avatar"] = avatarURL
	}

	// Define a struct for decoding the updated profile
	var updatedProfile models.User

	// Update the user's profile with the new data
	filter := bson.M{"_id": id}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = db.Collection(os.Getenv("USER_COLLECTION")).FindOneAndUpdate(c.Context(), filter, bson.M{"$set": updateData}, opts).Decode(&updatedProfile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error updating profile",
			"data":    err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Your profile's been updated!",
		"data":    updatedProfile,
	})
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
