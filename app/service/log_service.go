package service

import (
	"prak4/app/repository"
	"prak4/helper"
	"github.com/gofiber/fiber/v2"
)

type LogService struct {
	Repo *repository.LogRepository
}

// NewLogService creates a new instance of LogService
func NewLogService(repo *repository.LogRepository) *LogService {
	return &LogService{Repo: repo}
}

func (s *LogService) GetAllLogs(c *fiber.Ctx) error {
	logs, err := s.Repo.GetAllLogs()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data log")
	}
	return helper.SuccessResponse(c, logs, "Data log berhasil diambil")
}