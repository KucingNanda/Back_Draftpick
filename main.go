package main

import (
	"draftpick-api/config" // Mengambil fungsi dari folder config
	"draftpick-api/router"
	"log"
	"os"

	"draftpick-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	docs "draftpick-api/docs"
)

// @title Draft Pick API
// @version 1.0
// @description Dokumentasi API backend DraftPick menggunakan Golang Fiber, GORM, PostgreSQL, dan JWT.
// @contact.name Praktikum Pemrograman III
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Panggil fungsi koneksi ke database sebelum server menyala
	config.ConnectDB()

	// Inisialisasi Fiber dengan batas BodyLimit yang ditingkatkan (20 MB)
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024,
	})

	// Logger (Mencatat aktivitas API ke Terminal)
	app.Use(logger.New())

	// CORS
	middleware.SetupCORS(app)

	// Panggil Router
	router.SetupRoutes(app)

	// Setup Swagger Host secara dinamis dari ENV, Railway, atau Fallback
	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost == "" {
		swaggerHost = os.Getenv("RAILWAY_PUBLIC_DOMAIN")
	}
	if swaggerHost == "" {
		swaggerHost = "127.0.0.1:3000" // Fallback untuk lokal
	}
	docs.SwaggerInfo.Host = swaggerHost

	// Setup port dinamis untuk keperluan deployment (Alwaysdata/Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Port default lokal
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0" // Bind ke semua interface
	}

	// Jalankan server
	log.Fatal(app.Listen(host + ":" + port))
}
