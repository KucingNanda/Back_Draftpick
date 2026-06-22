package handler

import (
	"draftpick-api/config"
	"draftpick-api/model"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
// @Param request body map[string]string true "Payload register (username, password, role)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if input.Username == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username dan password wajib diisi"})
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengenkripsi password"})
	}

	role := "peserta"
	if input.Role == "wasit" || input.Role == "admin" {
		role = "wasit" // Normalisasi ke wasit
	}

	newUser := model.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
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

// Login godoc
// @Summary Login user
// @Description Melakukan autentikasi user dan mengembalikan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Payload login (username, password)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/login [post]
func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	var user model.User
	result := config.DB.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Username tidak ditemukan"})
	}

	// Verifikasi hash bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
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
// @Param request body map[string]string true "Payload ganti password (old_password, new_password)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/change-password [put]
func ChangePassword(c *fiber.Ctx) error {
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	userId := c.Locals("user_id") // Diambil dari token JWT

	var user model.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Verifikasi password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Password lama salah"})
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengenkripsi password baru"})
	}

	user.Password = string(hashedPassword)
	config.DB.Save(&user)

	return c.JSON(fiber.Map{"message": "Password berhasil diubah"})
}
