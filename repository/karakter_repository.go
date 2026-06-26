package repository

import (
	"draftpick-api/config"
	"draftpick-api/model"
)

// Mengambil semua data karakter
func GetAllKarakter() ([]model.Karakter, error) {
	var karakters []model.Karakter
	result := config.DB.Order("nama ASC").Find(&karakters)
	return karakters, result.Error
}

// Mengambil detail karakter berdasarkan ID
func GetKarakterByID(id string) (*model.Karakter, error) {
	var karakter model.Karakter
	result := config.DB.First(&karakter, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &karakter, nil
}

// Menyimpan karakter baru
func CreateKarakter(karakter *model.Karakter) error {
	result := config.DB.Create(karakter)
	return result.Error
}

// Menghapus karakter berdasarkan ID
func DeleteKarakter(id string) error {
	var karakter model.Karakter
	result := config.DB.Delete(&karakter, id)
	return result.Error
}

// Memperbarui data karakter
func UpdateKarakter(id string, karakter *model.Karakter) error {
	// Menggunakan map agar nilai 0 (seperti Cost) tetap di-update oleh GORM
	result := config.DB.Model(&model.Karakter{}).Where("id = ?", id).Updates(map[string]interface{}{
		"nama":        karakter.Nama,
		"elemen":      karakter.Elemen,
		"senjata":     karakter.Senjata,
		"rarity":      karakter.Rarity,
		"cost":        karakter.Cost,
		"image_url":   karakter.ImageURL,
		"konstelasi":  karakter.Konstelasi,
		"ulang_tahun": karakter.UlangTahun,
		"region":      karakter.Region,
		"sinergi":     karakter.Sinergi,
	})
	return result.Error
}
