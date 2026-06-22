package handler

import (
	"draftpick-api/config"
	"draftpick-api/model"

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
	var users []model.User
	// Ambil semua user dengan role peserta, tanpa menyertakan password
	result := config.DB.Select("id", "username", "role").Where("role = ?", "peserta").Find(&users)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data peserta"})
	}

	return c.JSON(users)
}
