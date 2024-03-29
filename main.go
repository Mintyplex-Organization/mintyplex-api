package main

import (
	"context"
	"fmt"
	"log"
	"mintyplex-api/internal/controllers"
	"mintyplex-api/internal/repository"
	"mintyplex-api/internal/routes"
	"mintyplex-api/internal/utils"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	// redisclient *redis.Client

	UserRepository      repository.UserRepository
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authRepo            repository.AuthRepository
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}
	fmt.Println("loaded..next")

	// uri := config.DBUri
	// if uri == "" {
	// 	log.Fatal("Set MONGODB_URI")
	// }

	// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	// if err != nil {
	// 	panic(err)
	// }
	// mongoclient = client

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// ctx = context.TODO()
	fmt.Println(ctx, "ctx reached")

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	fmt.Println(mongoclient, "mongo client resolved")

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	// redisclient = redis.NewClient(&redis.Options{
	// 	Addr: config.RedisUri,
	// })

	// if _, err := redisclient.Ping(ctx).Result(); err != nil {
	// 	panic(err)
	// }

	// err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Redis client connected successfully...")

	// Collections
	authCollection = mongoclient.Database("minty").Collection("users")
	UserRepository = repository.NewUserServiceImpl(authCollection, ctx)
	// authRepo = repository.NewAuthServiceImpl(authCollection, ctx)
	authRepo = repository.NewAuthRepository(authCollection, ctx)
	AuthController = controllers.NewAuthController(authRepo, UserRepository)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(UserRepository)
	UserRouteController = routes.NewRouteUserController(UserController)

	server = gin.Default()
}

func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}
	fmt.Println("loaded config..next")

	defer mongoclient.Disconnect(ctx)

	// value, err := redisclient.Get(ctx, "test").Result()

	// if err == redis.Nil {
	// 	fmt.Println("key: test does not exist")
	// } else if err != nil {
	// 	panic(err)
	// }

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "off to mars soon"})
	})

	AuthRouteController.AuthRoute(router, UserRepository)
	UserRouteController.UserRoute(router, UserRepository)
	log.Fatal(server.Run(":" + config.Port))
}
