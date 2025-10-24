// app/repository/user_mongo_repository.go
package repository

import (
	"context"
	"os"
	"prak4/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongo(db *mongo.Database) *UserRepositoryMongo {
	collectionName := os.Getenv("MONGO_COLLECTION_USERS")
	return &UserRepositoryMongo{
		collection: db.Collection(collectionName),
	}
}

func (r *UserRepositoryMongo) GetUserByUsername(ctx context.Context, username string) (*model.UserMongo, error) {
	var user model.UserMongo
	filter := bson.M{"username": username}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User tidak ditemukan, bukan error
		}
		return nil, err
	}
	return &user, nil
}