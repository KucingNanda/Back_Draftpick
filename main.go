package main

import (
	"draftpick-api/config" // Mengambil fungsi dari folder config
	"draftpick-api/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "draftpick-api/docs"
)

// @title Draft Pick API
// @version 1.0
// @description Dokumentasi API backend DraftPick menggunakan Golang Fiber, GORM, PostgreSQL, dan JWT.
// @contact.name Praktikum Pemrograman III
// @host 127.0.0.1:3000
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://front-draftpick.vercel.app", // Mengizinkan port lokal dan Vercel
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

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
