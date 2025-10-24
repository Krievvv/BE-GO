package service

import (
	"database/sql"
	"prak4/app/model"
	"prak4/app/repository"
	"prak4/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	UserRepo *repository.UserRepository
	LogRepo  *repository.LogRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo *repository.UserRepository, logRepo *repository.LogRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
		LogRepo:  logRepo,
	}
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

	token, err := utils.GenerateToken(*user,"postgres")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	logEntry := model.LogAktivitas{
		UserID:    user.ID,
		Username:  user.Username,
		Aksi:      "LOGIN_SUCCESS",
		Detail:    "Pengguna berhasil login.",
		Timestamp: time.Now(),
	}
	_ = s.LogRepo.CreateLog(logEntry)

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

func (s *AuthService) RegisterUser(c *fiber.Ctx) error {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// Validasi input
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, password, dan email harus diisi",
		})
	}

	// Hash password sebelum disimpan
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memproses password",
		})
	}

	// Set default role jika tidak diisi
	if user.Role == "" {
		user.Role = "user"
	}

	// Simpan user baru
	newUser, err := s.UserRepo.CreateUser(user.Username, user.Email, hashedPassword, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat user baru: " + err.Error(),
		})
	}

	// Log aktivitas registrasi
	logEntry := model.LogAktivitas{
		UserID:    newUser.ID,
		Username:  newUser.Username,
		Aksi:      "USER_REGISTERED",
		Detail:    "User baru berhasil didaftarkan",
		Timestamp: time.Now(),
	}
	_ = s.LogRepo.CreateLog(logEntry)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User berhasil didaftarkan",
		"data":    newUser,
	})
}