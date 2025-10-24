package service

import (
	"context"
	"prak4/app/model"
	repoMongo "prak4/app/repository" // Alias agar tidak konflik
	"prak4/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthMongoService struct {
	UserRepo *repoMongo.UserRepositoryMongo
	LogRepo  *repoMongo.LogRepository 
}

// Struct baru untuk respons login Mongo
type LoginMongoResponse struct {
	User  model.UserMongo `json:"user"`
	Token string          `json:"token"`
}

func NewAuthMongoService(userRepo *repoMongo.UserRepositoryMongo, logRepo *repoMongo.LogRepository) *AuthMongoService {
	return &AuthMongoService{UserRepo: userRepo, LogRepo: logRepo}
}

func (s *AuthMongoService) LoginMongo(c *fiber.Ctx) error {
	var req model.LoginRequest 
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username dan password harus diisi"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.UserRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error database"})
	}
	if user == nil || !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	token, err := utils.GenerateToken(model.User{ 
		ID:       user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	},"mongo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	logEntry := model.LogAktivitas{
		UserID:    user.UserID,
		Username:  user.Username,
		Aksi:      "LOGIN_MONGO_SUCCESS",
		Detail:    "Pengguna berhasil login via MongoDB auth.",
		Timestamp: time.Now(),
	}
	_ = s.LogRepo.CreateLog(logEntry)

	response := LoginMongoResponse{
		User: model.UserMongo{ 
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login MongoDB berhasil",
		"data":    response,
	})
}