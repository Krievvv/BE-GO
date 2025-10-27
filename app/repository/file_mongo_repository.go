package repository

import (
	"context"
	"errors"
	"prak4/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileRepositoryMongo struct {
	collection *mongo.Collection
}

func NewFileRepositoryMongo(db *mongo.Database) *FileRepositoryMongo {
	return &FileRepositoryMongo{
		collection: db.Collection("files"),
	}
}

func (r *FileRepositoryMongo) Create(ctx context.Context, file *model.File) (*model.File, error) {
	result, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		return nil, err
	}
	file.ID = result.InsertedID.(primitive.ObjectID)
	return file, nil
}

func (r *FileRepositoryMongo) FindAll(ctx context.Context) ([]model.File, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "uploaded_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []model.File
	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepositoryMongo) FindByID(ctx context.Context, id string) (*model.File, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID file tidak valid")
	}

	var file model.File
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil 
		}
		return nil, err
	}
	return &file, nil
}

func (r *FileRepositoryMongo) Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return r.collection.DeleteOne(ctx, bson.M{"_id": id})
}