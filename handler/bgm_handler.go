package handler

import (
	"draftpick-api/model"
	"draftpick-api/repository"

	"github.com/gofiber/fiber/v2"
)

// GetBgms godoc
// @Summary Ambil semua data BGM
// @Description Mengambil daftar seluruh BGM
// @Tags BGM
// @Produce json
// @Success 200 {array} model.Bgm
// @Failure 500 {object} map[string]string
// @Router /api/bgm [get]
func GetBgms(c *fiber.Ctx) error {
	bgms, err := repository.GetAllBgm()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data BGM"})
	}
	return c.JSON(bgms)
}

// GetBgmByID godoc
// @Summary Ambil detail data BGM berdasarkan ID
// @Description Mengambil satu data BGM spesifik berdasarkan ID
// @Tags BGM
// @Produce json
// @Param id path int true "ID BGM"
// @Success 200 {object} model.Bgm
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/bgm/{id} [get]
func GetBgmByID(c *fiber.Ctx) error {
	id := c.Params("id")
	bgm, err := repository.GetBgmByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "BGM tidak ditemukan"})
	}
	return c.JSON(bgm)
}

// AddBgm godoc
// @Summary Tambah data BGM
// @Description Menambahkan data BGM baru
// @Tags BGM
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Bgm true "Payload BGM"
// @Success 201 {object} model.Bgm
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/bgm [post]
func AddBgm(c *fiber.Ctx) error {
	var input model.Bgm
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	if input.Nama == "" || input.Url == "" {
		return c.Status(400).JSON(fiber.Map{"message": "nama dan url wajib diisi"})
	}

	if err := repository.CreateBgm(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan BGM"})
	}

	return c.Status(201).JSON(input)
}

// UpdateBgm godoc
// @Summary Ubah data BGM
// @Description Mengubah data BGM berdasarkan ID
// @Tags BGM
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID BGM"
// @Param request body model.Bgm true "Payload update"
// @Success 200 {object} model.Bgm
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/bgm/{id} [put]
func UpdateBgm(c *fiber.Ctx) error {
	id := c.Params("id")
	bgm, err := repository.GetBgmByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "BGM tidak ditemukan"})
	}

	var input model.Bgm
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	if input.Nama == "" {
		return c.Status(400).JSON(fiber.Map{"message": "nama wajib diisi"})
	}

	bgm.Nama = input.Nama
	if input.Url != "" {
		bgm.Url = input.Url
	}

	if err := repository.SaveBgm(&bgm); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengupdate BGM"})
	}

	return c.JSON(bgm)
}

// DeleteBgm godoc
// @Summary Hapus data BGM
// @Description Menghapus data BGM berdasarkan ID
// @Tags BGM
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID BGM"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/bgm/{id} [delete]
func DeleteBgm(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repository.DeleteBgm(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus BGM"})
	}

	return c.JSON(fiber.Map{"message": "BGM berhasil dihapus"})
}
