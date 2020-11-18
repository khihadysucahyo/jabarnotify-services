package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB init dbhost
var MongoDB *mongo.Database

// InitMongoDB to init dbhost
func InitMongoDB(url string, nameDB string) *mongo.Database {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(url),
	)

	if err != nil {
		log.Fatal(err)
	}
	MongoDB = client.Database(nameDB)
	return MongoDB
}

// GetMongoDB Using this function to get a connection, you can create your connection pool here.
func GetMongoDB() *mongo.Database {
	return MongoDB
}

// ConfigDB mongo
type ConfigDB struct {
	DB *mongo.Database
}

// Initialize mongodb connection
func Initialize() (*ConfigDB, error) {
	addr := "mongodb://" + GetEnv("DB_MONGO_HOST") + ":" + GetEnv("DB_MONGO_PORT")
	nameDB := GetEnv("DB_MONGO_NAME")
	config := ConfigDB{}
	// Connect to MongoDB
	ForeverSleep(2*time.Second, func(attempt int) error {
		db := InitMongoDB(addr, nameDB)
		config.DB = db
		return nil
	})
	return &config, nil
}
