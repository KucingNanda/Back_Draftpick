package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected melindungi rute dan memastikan JWT dikirimkan dan valid
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Token tidak ditemukan atau format salah",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "rahasia-draftpick-default" // Fallback jika .env belum diset
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Token tidak valid atau sudah kedaluwarsa",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Klaim token tidak valid",
			})
		}

		// Simpan claims ke locals agar bisa diakses oleh handler
		c.Locals("user_id", claims["id"])
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

// RequireRole membatasi akses hanya untuk role tertentu (misal: "wasit")
func RequireRole(allowedRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != allowedRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
			})
		}
		return c.Next()
	}
}
