package handler

import (
	"draftpick-api/model"
	"draftpick-api/repository"
	"draftpick-api/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT membuat token JWT untuk user
func GenerateJWT(user model.User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "rahasia-draftpick-default"
	}

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 72 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// Register godoc
// @Summary Register user baru
// @Description Mendaftarkan peserta atau admin baru ke dalam sistem
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.AuthRequest true "Payload register (username, password)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {
	var input model.AuthRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if input.Username == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username dan password wajib diisi"})
	}

	// Hash password menggunakan utils
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengenkripsi password"})
	}

	role := "member"

	newUser := model.User{
		Username: input.Username,
		Password: hashedPassword,
		Role:     role,
	}

	if err := repository.CreateUser(&newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username mungkin sudah digunakan"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Pendaftaran berhasil",
		"user": fiber.Map{
			"id":       newUser.ID,
			"username": newUser.Username,
			"role":     newUser.Role,
		},
	})
}

// GetMe godoc
// @Summary Ambil profil user yang sedang login
// @Description Mengembalikan informasi user saat ini beserta statusnya
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/me [get]
func GetMe(c *fiber.Ctx) error {
	userId := c.Locals("user_id")
	user, err := repository.GetUserByID(userId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	peserta, _ := repository.GetPesertaByUserID(user.ID)

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"peserta":  peserta,
	})
}


// Login godoc
// @Summary Login user
// @Description Melakukan autentikasi user dan mengembalikan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Payload login (username, password)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/login [post]
func Login(c *fiber.Ctx) error {
	var input model.LoginRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	user, err := repository.GetUserByUsername(input.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Username tidak ditemukan"})
	}

	// Verifikasi hash
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
	}

	// Generate JWT
	token, err := GenerateJWT(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghasilkan token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
		"role":    user.Role,
	})
}

// ChangePassword godoc
// @Summary Ubah password user
// @Description Mengubah password user yang sedang login
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordRequest true "Payload ganti password (old_password, new_password)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/change-password [put]
func ChangePassword(c *fiber.Ctx) error {
	var input model.ChangePasswordRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	userId := c.Locals("user_id") // Diambil dari token JWT

	user, err := repository.GetUserByID(userId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Verifikasi password lama
	if !utils.CheckPasswordHash(input.OldPassword, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Password lama salah"})
	}

	// Hash password baru
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengenkripsi password baru"})
	}

	user.Password = hashedPassword
	repository.SaveUser(&user)

	return c.JSON(fiber.Map{"message": "Password berhasil diubah"})
}
