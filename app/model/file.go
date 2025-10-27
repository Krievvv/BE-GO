package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// metadata file di MongoDB
type File struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       int                `json:"user_id" bson:"user_id"` 
	FileName     string             `json:"file_name" bson:"file_name"` 
	OriginalName string             `json:"original_name" bson:"original_name"` 
	FilePath     string             `json:"file_path" bson:"file_path"` 
	FileSize     int64              `json:"file_size" bson:"file_size"` 
	FileType     string             `json:"file_type" bson:"file_type"` 
	UploadedAt   time.Time          `json:"uploaded_at" bson:"uploaded_at"`
	Jenis        string             `json:"jenis" bson:"jenis"` 
}

// struct untuk respons API (ID sebagai string)
type FileResponse struct {
	ID           string    `json:"id"`
	UserID       int       `json:"user_id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	FilePath     string    `json:"file_path"` // berisi URL untuk akses
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	Jenis        string    `json:"jenis"`
}