package service

import (
	"database/sql"
	"prak4/app/model"
	"prak4/app/repository"
	"prak4/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username dan password harus diisi"})
	}

	user, passwordHash, err := s.UserRepo.GetUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau password salah"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error database"})
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	token, err := utils.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	response := model.LoginResponse{
		User:  *user,
		Token: token,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}