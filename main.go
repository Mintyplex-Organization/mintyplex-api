package main

import (
	"context"
	"fmt"
	"log"
	"mintyplex-api/internal/controllers"
	"mintyplex-api/internal/database"
	"mintyplex-api/internal/middleware"
	"mintyplex-api/internal/routes"
	"mintyplex-api/internal/services"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	server      *fiber.App
	ctx         context.Context
	mongoclient *mongo.Client
	db          *mongo.Database

	userService         services.UserService
	userController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		panic("Error Loading .env File, Check If It Exists.")
	}
	log.Println("Connecting to MongoDB...")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_SRV_RECORD")).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err, "this is me")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	mongoclient = client
	db = client.Database(os.Getenv("MONGODB_DATABASE"))
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 200 * 1024 * 1024, // this is the default limit of 4MB
	})

	corsConfig := cors.Config{
		AllowOrigins:     "http://localhost:8000, http://localhost:3000",
		AllowCredentials: true,
	}
	app.Use(cors.New(corsConfig))

	middleware.CorsMiddleware(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Plexer SZN")
	})

	db := database.MongoClient()
	app.Use(middleware.IngestDb(db))

	AuthRouteController.AuthRoute(app, userService)
	UserRouteController.UserRoutes(app, userService)


	log.Fatal(app.Listen(":8081"))

	// app.Listen(":8081")
}
