package helper

import "github.com/gofiber/fiber/v2"

// SuccessResponse membuat respons JSON untuk kasus sukses
func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse membuat respons JSON untuk kasus error
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}