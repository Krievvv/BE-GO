package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Database

func ConnectMongo() {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE")

	if uri == "" || dbName == "" {
		log.Fatal("MONGO_URI dan MONGO_DATABASE harus diatur di file .env")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Gagal membuat client MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}
	
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Gagal ping MongoDB: %v", err)
	}

	log.Println("Berhasil terhubung ke MongoDB")
	Mongo = client.Database(dbName)
}