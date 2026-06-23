package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS mengkonfigurasi aturan akses lintas domain (CORS)
func SetupCORS(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://front-draftpick.vercel.app", // Mengizinkan port lokal dan Vercel
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}
