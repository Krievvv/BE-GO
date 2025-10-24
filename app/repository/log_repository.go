package repository

import (
	"context"
	"prak4/app/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository struct {
	DB *mongo.Database
}

// NewLogRepository creates a new instance of LogRepository
func NewLogRepository(db *mongo.Database) *LogRepository {
	return &LogRepository{DB: db}
}

// CreateLog menyimpan log baru ke MongoDB
func (r *LogRepository) CreateLog(log model.LogAktivitas) error {
	_, err := r.DB.Collection("aktivitas").InsertOne(context.Background(), log)
	return err
}

// GetAllLogs mengambil semua log dari MongoDB
func (r *LogRepository) GetAllLogs() ([]model.LogAktivitas, error) {
	cursor, err := r.DB.Collection("aktivitas").Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var logs []model.LogAktivitas
	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}
	return logs, nil
}