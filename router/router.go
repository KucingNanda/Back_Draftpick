package router

import (
	"draftpick-api/handler"
	"draftpick-api/middleware"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Rute Publik
	app.Get("/", handler.GetRoot)
	app.Get("/docs/*", swagger.HandlerDefault)
	api.Get("/status", handler.GetStatus)
	
	// API Publik (Hanya GET)
	api.Get("/karakter", handler.GetKarakters)
	api.Get("/karakter/:id", handler.GetKarakterByID)
	api.Get("/bgm", handler.GetBgms)
	api.Get("/bgm/:id", handler.GetBgmByID)
	api.Get("/settings", handler.GetSettings)

	// Rute Auth
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)

	// Rute Proteksi (Hanya yang punya token JWT valid)
	protected := api.Group("/", middleware.Protected())

	// Users (Protected)
	protected.Get("/me", handler.GetMe)
	protected.Get("/users", middleware.RequireRole("wasit"), handler.GetUsers)
	protected.Put("/users/:id/role", middleware.RequireRole("wasit"), handler.UpdateUserRole)

	// Karakter (Protected)
	protected.Post("/karakter", middleware.RequireRole("wasit"), handler.AddKarakter)          
	protected.Put("/karakter/:id", middleware.RequireRole("wasit"), handler.UpdateKarakter)    
	protected.Delete("/karakter/:id", middleware.RequireRole("wasit"), handler.DeleteKarakter) 

	// Peserta (Protected)
	protected.Get("/peserta", handler.GetPesertas)
	protected.Get("/peserta/:id", handler.GetPesertaByID)
	protected.Get("/my-peserta", handler.GetMyPeserta)
	protected.Post("/my-peserta", handler.SaveMyPeserta)
	protected.Post("/peserta", middleware.RequireRole("wasit"), handler.AddPeserta)
	protected.Put("/peserta/:id", middleware.RequireRole("wasit"), handler.UpdatePeserta)
	protected.Delete("/peserta/:id", middleware.RequireRole("wasit"), handler.DeletePeserta)

	// BGM (Protected)
	protected.Post("/bgm", middleware.RequireRole("wasit"), handler.AddBgm)
	protected.Put("/bgm/:id", middleware.RequireRole("wasit"), handler.UpdateBgm)
	protected.Delete("/bgm/:id", middleware.RequireRole("wasit"), handler.DeleteBgm)

	// History (Protected)
	protected.Post("/history", handler.SaveHistory) // Semua bisa save (atau biarkan jika draft bisa dimainkan peserta)
	protected.Get("/history", handler.GetHistories)
	protected.Delete("/history/:id", middleware.RequireRole("wasit"), handler.DeleteHistory)

	// Rute Upload Cloudinary (Protected)
	protected.Post("/upload", middleware.RequireRole("wasit"), handler.UploadImage)
	protected.Post("/upload/audio", middleware.RequireRole("wasit"), handler.UploadAudio)

	// Settings (Protected)
	protected.Put("/settings", middleware.RequireRole("wasit"), handler.UpdateSettings)


	// Ganti password (semua user login)
	protected.Put("/change-password", handler.ChangePassword)
}
