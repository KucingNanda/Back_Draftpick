package handler

import (
	"draftpick-api/config"
	"draftpick-api/model"

	"github.com/gofiber/fiber/v2"
)

type HistoryRequest struct {
	Role      string `json:"role"`
	PlayerA   string `json:"playerA"`
	PlayerB   string `json:"playerB"`
	DraftData string `json:"draftData"`
}

// SaveHistory godoc
// @Summary Simpan riwayat pertandingan
// @Description Menyimpan data history match
// @Tags History
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.History true "Payload history"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/history [post]
func SaveHistory(c *fiber.Ctx) error {
	var req HistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if req.Role != "wasit" && req.Role != "peserta" {
		return c.Status(403).JSON(fiber.Map{"error": "Role tidak diizinkan menyimpan riwayat"})
	}

	limit := int64(10)
	if req.Role == "wasit" {
		limit = 20
	}

	var count int64
	config.DB.Model(&model.History{}).Where("role = ?", req.Role).Count(&count)

	if count >= limit {
		// Limit tercapai, timpa riwayat yang paling lama (UpdatedAt paling tua)
		var oldest model.History
		if err := config.DB.Where("role = ?", req.Role).Order("updated_at asc").First(&oldest).Error; err == nil {
			oldest.PlayerA = req.PlayerA
			oldest.PlayerB = req.PlayerB
			oldest.DraftData = req.DraftData
			// GORM otomatis memperbarui UpdatedAt
			config.DB.Save(&oldest)
			return c.JSON(fiber.Map{"message": "Riwayat lama berhasil ditimpa", "history": oldest})
		}
	}

	// Buat baru jika belum mencapai limit
	newHistory := model.History{
		Role:      req.Role,
		PlayerA:   req.PlayerA,
		PlayerB:   req.PlayerB,
		DraftData: req.DraftData,
	}

	if err := config.DB.Create(&newHistory).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan riwayat"})
	}

	return c.JSON(fiber.Map{"message": "Riwayat berhasil disimpan", "history": newHistory})
}

// GetHistories godoc
// @Summary Ambil riwayat pertandingan
// @Description Mengambil daftar seluruh history match
// @Tags History
// @Produce json
// @Success 200 {array} model.History
// @Failure 500 {object} map[string]string
// @Router /api/history [get]
func GetHistories(c *fiber.Ctx) error {
	role := c.Query("role")
	if role == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Parameter role diperlukan"})
	}

	var histories []model.History
	// Mengurutkan berdasarkan yang terbaru diedit/ditimpa/dibuat
	config.DB.Where("role = ?", role).Order("updated_at desc").Find(&histories)

	return c.JSON(histories)
}

// DeleteHistory godoc
// @Summary Hapus riwayat pertandingan
// @Description Menghapus history match berdasarkan ID
// @Tags History
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID History"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/history/{id} [delete]
func DeleteHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&model.History{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus riwayat"})
	}
	return c.JSON(fiber.Map{"message": "Riwayat berhasil dihapus"})
}
