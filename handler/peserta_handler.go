package handler

import (
	"draftpick-api/model"
	"draftpick-api/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Ambil Semua Peserta
// GetPesertas godoc
// @Summary Ambil semua data peserta
// @Description Mengambil daftar seluruh peserta
// @Tags Peserta
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Peserta
// @Failure 500 {object} map[string]string
// @Router /api/peserta [get]
func GetPesertas(c *fiber.Ctx) error {
	pesertas, err := repository.GetAllPeserta()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data peserta"})
	}
	return c.JSON(pesertas)
}

// GetPesertaByID godoc
// @Summary Ambil detail data peserta berdasarkan ID
// @Description Mengambil satu data peserta spesifik berdasarkan ID
// @Tags Peserta
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Peserta"
// @Success 200 {object} model.Peserta
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/peserta/{id} [get]
func GetPesertaByID(c *fiber.Ctx) error {
	id := c.Params("id")
	peserta, err := repository.GetPesertaByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Peserta tidak ditemukan"})
	}
	return c.JSON(peserta)
}

// Tambah Peserta Baru
// AddPeserta godoc
// @Summary Tambah data peserta
// @Description Menambahkan data peserta baru
// @Tags Peserta
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Peserta true "Payload peserta"
// @Success 201 {object} model.Peserta
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/peserta [post]
func AddPeserta(c *fiber.Ctx) error {
	peserta := new(model.Peserta)
	if err := c.BodyParser(peserta); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input tidak valid"})
	}

	if peserta.Nama == "" {
		return c.Status(400).JSON(fiber.Map{"message": "nama wajib diisi"})
	}

	if len(peserta.Karakters) != 8 {
		return c.Status(400).JSON(fiber.Map{"error": "Peserta harus mendaftarkan tepat 8 Karakter Limited"})
	}

	err := repository.CreatePeserta(peserta)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan peserta"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Peserta berhasil ditambahkan!",
		"data":    peserta,
	})
}

// Update Peserta
// UpdatePeserta godoc
// @Summary Ubah data peserta
// @Description Mengubah data peserta berdasarkan ID
// @Tags Peserta
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Peserta"
// @Param request body model.Peserta true "Payload update"
// @Success 200 {object} model.Peserta
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/peserta/{id} [put]
func UpdatePeserta(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	peserta := new(model.Peserta)
	if err := c.BodyParser(peserta); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input tidak valid"})
	}
	peserta.ID = uint(id)

	if peserta.Nama == "" {
		return c.Status(400).JSON(fiber.Map{"message": "nama wajib diisi"})
	}

	if len(peserta.Karakters) != 8 {
		return c.Status(400).JSON(fiber.Map{"error": "Peserta harus mendaftarkan tepat 8 Karakter Limited"})
	}

	err = repository.UpdatePeserta(idStr, peserta)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memperbarui peserta"})
	}

	return c.JSON(fiber.Map{
		"message": "Peserta berhasil diperbarui!",
	})
}

// Hapus Peserta
// DeletePeserta godoc
// @Summary Hapus data peserta
// @Description Menghapus data peserta berdasarkan ID
// @Tags Peserta
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Peserta"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/peserta/{id} [delete]
func DeletePeserta(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.DeletePeserta(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus peserta"})
	}

	return c.JSON(fiber.Map{"message": "Peserta berhasil dihapus"})
}

// GetMyPeserta godoc
// @Summary Ambil biodata & karakter peserta milik akun sendiri
// @Description Mengambil data profil peserta yang tertaut dengan akun user yang sedang login
// @Tags Peserta
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.Peserta
// @Failure 404 {object} map[string]string
// @Router /api/my-peserta [get]
func GetMyPeserta(c *fiber.Ctx) error {
	userId := c.Locals("user_id")
	peserta, err := repository.GetPesertaByUserID(userId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data profil"})
	}
	if peserta == nil {
		// Mengembalikan 200 OK dengan null jika user belum pernah membuat formasi
		return c.JSON(nil)
	}
	return c.JSON(peserta)
}

// SaveMyPeserta godoc
// @Summary Simpan atau perbarui biodata & 8 karakter milik akun sendiri
// @Description Menyimpan formasi 8 karakter Limited untuk akun member/peserta yang sedang login
// @Tags Peserta
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Peserta true "Payload peserta"
// @Success 200 {object} model.Peserta
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/my-peserta [post]
func SaveMyPeserta(c *fiber.Ctx) error {
	userId := c.Locals("user_id")
	var uID uint
	switch v := userId.(type) {
	case float64:
		uID = uint(v)
	case int:
		uID = uint(v)
	case uint:
		uID = v
	case string:
		idParsed, _ := strconv.ParseUint(v, 10, 32)
		uID = uint(idParsed)
	}

	if uID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	peserta := new(model.Peserta)
	if err := c.BodyParser(peserta); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input tidak valid"})
	}

	if peserta.Nama == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Nama peserta wajib diisi"})
	}

	if len(peserta.Karakters) != 8 {
		return c.Status(400).JSON(fiber.Map{"error": "Peserta harus mendaftarkan tepat 8 Karakter Limited"})
	}

	savedPeserta, err := repository.SaveMyPeserta(uID, peserta)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan formasi peserta"})
	}

	return c.JSON(fiber.Map{
		"message": "Profil peserta berhasil disimpan!",
		"data":    savedPeserta,
	})
}

