package repository

import (
	"draftpick-api/config"
	"draftpick-api/model"
)

// Mengambil semua data peserta beserta detail karakter Limited 5-Star miliknya
func GetAllPeserta() ([]model.Peserta, error) {
	var pesertas []model.Peserta
	// Preload "Karakters" (tabel join) dan "Karakters.Karakter" (detail karakternya)
	result := config.DB.Preload("Karakters").Preload("Karakters.Karakter").Order("nama ASC").Find(&pesertas)
	return pesertas, result.Error
}

// Mengambil detail peserta berdasarkan ID
func GetPesertaByID(id string) (*model.Peserta, error) {
	var peserta model.Peserta
	result := config.DB.Preload("Karakters").Preload("Karakters.Karakter").First(&peserta, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &peserta, nil
}

// Menyimpan peserta baru beserta susunan 8 karakternya
func CreatePeserta(peserta *model.Peserta) error {
	result := config.DB.Create(peserta)
	return result.Error
}

// Menghapus peserta berdasarkan ID
func DeletePeserta(id string) error {
	var peserta model.Peserta
	// Penghapusan tabel PesertaKarakter akan diurus otomatis oleh GORM berkat OnDelete:CASCADE
	result := config.DB.Delete(&peserta, id)
	return result.Error
}

// Memperbarui data peserta (termasuk mengganti formasi 8 karakternya)
func UpdatePeserta(id string, peserta *model.Peserta) error {
	// Mulai Transaksi Database
	tx := config.DB.Begin()

	// 1. Update Nama Peserta
	if err := tx.Model(&model.Peserta{}).Where("id = ?", id).Update("nama", peserta.Nama).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Hapus semua data PesertaKarakter lama
	if err := tx.Where("peserta_id = ?", id).Delete(&model.PesertaKarakter{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. Masukkan data PesertaKarakter baru
	for i := range peserta.Karakters {
		// Set ulang PesertaID agar sesuai dengan rute update
		peserta.Karakters[i].PesertaID = peserta.ID
	}
	if err := tx.Create(&peserta.Karakters).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Selesaikan Transaksi
	return tx.Commit().Error
}
