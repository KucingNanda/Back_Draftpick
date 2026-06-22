package handler

import "github.com/gofiber/fiber/v2"

// GetRoot menangani rute utama (Root URL)
func GetRoot(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Selamat Datang di Server API Draft Pick 2.0!",
		"status":  "Server Aktif",
	})
}

// GetStatus menangani rute untuk pengecekan status API
func GetStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Halo! Mesin API Draft Pick Versi 2.0 sudah berjalan lancar",
	})
}
