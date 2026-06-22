package handler

import (
	"draftpick-api/model"
	"draftpick-api/repository"

	"github.com/gofiber/fiber/v2"
)

// Menghitung Cost berdasarkan Rarity
func calculateCost(rarity string) int {
	if rarity == "5 Star Limited" {
		return 1
	}
	return 0
}

// Handler untuk melihat semua karakter
// GetKarakters godoc
// @Summary Ambil semua data karakter
// @Description Mengambil daftar seluruh karakter
// @Tags Karakter
// @Produce json
// @Success 200 {array} model.Karakter
// @Failure 500 {object} map[string]string
// @Router /api/karakter [get]
func GetKarakters(c *fiber.Ctx) error {
	karakters, err := repository.GetAllKarakter()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data karakter"})
	}
	return c.JSON(karakters)
}

// GetKarakterByID godoc
// @Summary Ambil detail data karakter berdasarkan ID
// @Description Mengambil satu data karakter spesifik berdasarkan ID
// @Tags Karakter
// @Produce json
// @Param id path int true "ID Karakter"
// @Success 200 {object} model.Karakter
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/karakter/{id} [get]
func GetKarakterByID(c *fiber.Ctx) error {
	id := c.Params("id")
	karakter, err := repository.GetKarakterByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Karakter tidak ditemukan"})
	}
	return c.JSON(karakter)
}

// Handler untuk menambah karakter baru
// AddKarakter godoc
// @Summary Tambah data karakter
// @Description Menambahkan data karakter baru (hanya wasit)
// @Tags Karakter
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Karakter true "Payload karakter"
// @Success 201 {object} model.Karakter
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/karakter [post]
func AddKarakter(c *fiber.Ctx) error {
	karakter := new(model.Karakter)

	if err := c.BodyParser(karakter); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input tidak valid"})
	}

	if karakter.Nama == "" {
		return c.Status(400).JSON(fiber.Map{"message": "nama wajib diisi"})
	}
	if karakter.Rarity == "" {
		return c.Status(400).JSON(fiber.Map{"message": "rarity wajib diisi"})
	}

	// Auto-calculate cost
	karakter.Cost = calculateCost(karakter.Rarity)

	err := repository.CreateKarakter(karakter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan karakter"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Karakter berhasil ditambahkan!",
		"data":    karakter,
	})
}

// Handler untuk memperbarui karakter
// UpdateKarakter godoc
// @Summary Ubah data karakter
// @Description Mengubah data karakter berdasarkan ID (hanya wasit)
// @Tags Karakter
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Karakter"
// @Param request body model.Karakter true "Payload update"
// @Success 200 {object} model.Karakter
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/karakter/{id} [put]
func UpdateKarakter(c *fiber.Ctx) error {
	id := c.Params("id")
	karakter := new(model.Karakter)

	if err := c.BodyParser(karakter); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input tidak valid"})
	}

	if karakter.Nama == "" {
		return c.Status(400).JSON(fiber.Map{"message": "nama wajib diisi"})
	}
	if karakter.Rarity == "" {
		return c.Status(400).JSON(fiber.Map{"message": "rarity wajib diisi"})
	}

	// Auto-calculate cost just in case rarity was updated
	karakter.Cost = calculateCost(karakter.Rarity)

	err := repository.UpdateKarakter(id, karakter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memperbarui karakter"})
	}

	return c.JSON(fiber.Map{
		"message": "Karakter berhasil diperbarui!",
		"data":    karakter,
	})
}

// Handler untuk menghapus karakter
// DeleteKarakter godoc
// @Summary Hapus data karakter
// @Description Menghapus data karakter berdasarkan ID (hanya wasit)
// @Tags Karakter
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Karakter"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/karakter/{id} [delete]
func DeleteKarakter(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.DeleteKarakter(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus karakter"})
	}

	return c.JSON(fiber.Map{"message": "Karakter berhasil dihapus"})
}
