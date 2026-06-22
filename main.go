package main

import (
	"draftpick-api/config" // Mengambil fungsi dari folder config
	"draftpick-api/router"
	"log"

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
		AllowOrigins: "http://localhost:5173", // Mengizinkan port Vite React
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Panggil Router
	router.SetupRoutes(app)

	// Jalankan server di port 3000
	log.Fatal(app.Listen(":3000"))
}
