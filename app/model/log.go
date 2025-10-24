package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogAktivitas struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    int                `json:"user_id" bson:"user_id"`
	Username  string             `json:"username" bson:"username"`
	Aksi      string             `json:"aksi" bson:"aksi"`
	Detail    string             `json:"detail" bson:"detail"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}