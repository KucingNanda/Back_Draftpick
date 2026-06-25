package handler

import (
	"draftpick-api/model"
	"draftpick-api/repository"

	"github.com/gofiber/fiber/v2"
)

// GetSettings mengembalikan pengaturan saat ini
// GetSettings godoc
// @Summary Ambil pengaturan turnamen
// @Description Mengambil konfigurasi draft pick turnamen
// @Tags Settings
// @Produce json
// @Success 200 {object} model.Settings
// @Failure 500 {object} map[string]string
// @Router /api/settings [get]
func GetSettings(c *fiber.Ctx) error {
	settings, err := repository.GetSettings()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil pengaturan"})
	}

	return c.JSON(settings)
}

// UpdateSettings mengubah pengaturan (dilindungi di frontend via admin)
// UpdateSettings godoc
// @Summary Ubah pengaturan turnamen
// @Description Memperbarui konfigurasi max cost, bgm, dan special character
// @Tags Settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Settings true "Payload settings"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/settings [put]
func UpdateSettings(c *fiber.Ctx) error {
	var input model.Settings
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data pengaturan tidak valid"})
	}

	settings, err := repository.GetSettings()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menemukan pengaturan"})
	}

	settings.MaxCost = input.MaxCost
	settings.Free4StarId = input.Free4StarId
	
	if input.Free4StarHalf == "" {
		settings.Free4StarHalf = "all"
	} else {
		settings.Free4StarHalf = input.Free4StarHalf
	}
	
	settings.BgmUrl = input.BgmUrl

	if err := repository.SaveSettings(&settings); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan pengaturan"})
	}

	return c.JSON(fiber.Map{
		"message": "Pengaturan berhasil diperbarui",
		"settings": settings,
	})
}
