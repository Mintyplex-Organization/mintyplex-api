package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient() *mongo.Database {
	err := godotenv.Load()
	if err != nil {
		panic("Error Loading .env File, Check If It Exists.")
	}
	log.Println("Connecting to MongoDB...")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_SRV_RECORD")).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println(os.Getenv("MONGODB_SRV_RECORD"))

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

	return client.Database(os.Getenv("MONGODB_DATABASE"))
}
