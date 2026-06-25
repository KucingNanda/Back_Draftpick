package handler

import (
	"draftpick-api/repository"

	"github.com/gofiber/fiber/v2"
)

// GetUsers godoc
// @Summary Ambil daftar akun pengguna (Peserta)
// @Description Mengambil semua data pengguna dengan role peserta
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]string
// @Router /api/users [get]
func GetUsers(c *fiber.Ctx) error {
	users, err := repository.GetUsersByRole("peserta")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data peserta"})
	}

	return c.JSON(users)
}
