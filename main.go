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

	// Inisialisasi Fiber
	app := fiber.New()

	// Logger (Mencatat aktivitas API ke Terminal)
	app.Use(logger.New())

	// CORS
	middleware.SetupCORS(app)

	// Panggil Router
	router.SetupRoutes(app)

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
