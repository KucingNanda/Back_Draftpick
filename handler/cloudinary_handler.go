package handler

import (
	"context"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
)

// Handler untuk mengunggah gambar ke Cloudinary
func UploadImage(c *fiber.Ctx) error {
	// Ambil file dari form input (name="image")
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal mengambil file gambar dari request"})
	}

	// Buka file yang diupload
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca file gambar"})
	}
	defer file.Close()

	// Inisialisasi Cloudinary menggunakan CLOUDINARY_URL dari env
	// Pastikan CLOUDINARY_URL terisi di .env
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kredensial Cloudinary tidak valid atau tidak ditemukan"})
	}

	// Konfigurasi Upload
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "Draftpick/Karakter",
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengunggah gambar ke Cloudinary"})
	}

	// Berhasil
	return c.JSON(fiber.Map{
		"message":  "Gambar berhasil diunggah",
		"imageUrl": uploadResult.SecureURL,
	})
}

// Handler untuk mengunggah file Audio (BGM) ke Cloudinary
func UploadAudio(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("audio")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal mengambil file audio dari request"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca file audio"})
	}
	defer file.Close()

	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kredensial Cloudinary tidak valid"})
	}

	// Gunakan context dengan timeout (misal 60 detik) untuk mencegah proses menggantung selamanya
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       "Draftpick/Assets",
		ResourceType: "auto", // Gunakan auto agar Cloudinary mendeteksi format yang tepat, jangan paksakan "video"
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengunggah audio ke Cloudinary: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":  "Audio berhasil diunggah",
		"audioUrl": uploadResult.SecureURL,
	})
}
