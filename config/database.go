package config

import (
	"fmt"
	"log"
	"os"

	"draftpick-api/model" // Import folder model kita
	"draftpick-api/utils"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		_ = godotenv.Load("BackEnd/.env")
		_ = godotenv.Load("../BackEnd/.env")
	}

	dsn := os.Getenv("SUPABASE_DSN")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
	}

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // WAJIB untuk Supabase connection pooler (port 6543)
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatal("❌ Gagal terhubung ke database PostgreSQL:\n", err)
	}

	DB = database
	fmt.Println("✅ Koneksi ke PostgreSQL BERHASIL!")

	// Menjalankan Auto Migration untuk membuat tabel otomatis
	err = DB.AutoMigrate(&model.Karakter{}, &model.User{}, &model.Peserta{}, &model.PesertaKarakter{}, &model.Settings{}, &model.Bgm{}, &model.History{})
	if err != nil {
		log.Fatal("❌ Gagal melakukan migrasi database:\n", err)
	}

	// Buat pengaturan default jika belum ada
	var count int64
	DB.Model(&model.Settings{}).Count(&count)
	if count == 0 {
		DB.Create(&model.Settings{
			MaxCost:     7,
			Free4StarId: 0,
			BgmUrl:      "",
		})
	}
	// Hash password menggunakan utils
	adminHash, _ := utils.HashPassword("kucing0909")
	pesertaHash, _ := utils.HashPassword("peserta123")

	// Setup kredensial admin
	var adminUser model.User
	DB.Where("username = ?", "kucingabu").Find(&adminUser)
	if adminUser.ID == 0 {
		DB.Create(&model.User{
			Username: "kucingabu",
			Password: adminHash,
			Role:     "wasit",
		})
	}

	// Setup kredensial peserta latihan
	var pesertaUser model.User
	DB.Where("username = ?", "peserta").Find(&pesertaUser)
	if pesertaUser.ID == 0 {
		DB.Create(&model.User{
			Username: "peserta",
			Password: pesertaHash,
			Role:     "peserta",
		})
	}

	fmt.Println("✅ Migrasi Tabel Database Selesai!")
}
