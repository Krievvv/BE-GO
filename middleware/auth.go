package middleware

import (
	"prak4/utils"
	"strings"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token akses diperlukan"})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Format token tidak valid"})
		}
		
		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau kedaluwarsa"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("issuer", claims.Issuer)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan."})
		}
		return c.Next()
	}
}

func RequireIssuer(expectedIssuer string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		issuer, ok := c.Locals("issuer").(string)
		if !ok || issuer != expectedIssuer {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Akses Ditolak",
				"message": "Token tidak valid untuk sistem ini.",
			})
		}
		return c.Next()
	}
}