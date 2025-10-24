package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserMongo struct {
	ID           primitive.ObjectID `json:"-" bson:"_id,omitempty"` 
	UserID       int                `json:"user_id" bson:"user_id"` 
	Username     string             `json:"username" bson:"username"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"` 
	Role         string             `json:"role" bson:"role"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}