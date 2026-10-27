package repository

import (
	"draftpick-api/config"
	"draftpick-api/model"
	"fmt"
)

// Mengambil semua data peserta resmi (peserta manual wasit atau member yang sudah dipromosikan jadi peserta)
func GetAllPeserta() ([]model.Peserta, error) {
	var pesertas []model.Peserta
	result := config.DB.
		Joins("LEFT JOIN users ON users.id = peserta.user_id").
		Where("peserta.user_id IS NULL OR users.role = ?", "peserta").
		Preload("Karakters").
		Preload("Karakters.Karakter").
		Order("peserta.nama ASC").
		Find(&pesertas)
	return pesertas, result.Error
}

// Mengambil detail peserta berdasarkan ID
func GetPesertaByID(id string) (*model.Peserta, error) {
	var peserta model.Peserta
	result := config.DB.Preload("User").Preload("Karakters").Preload("Karakters.Karakter").First(&peserta, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &peserta, nil
}

// Mengambil profil peserta milik akun user tertentu
func GetPesertaByUserID(userID interface{}) (*model.Peserta, error) {
	var pesertas []model.Peserta
	err := config.DB.Preload("Karakters").Preload("Karakters.Karakter").Where("user_id = ?", userID).Limit(1).Find(&pesertas).Error
	if err != nil {
		return nil, err
	}
	if len(pesertas) == 0 {
		return nil, nil // Belum ada formasi yang dibuat (bukan error DB)
	}
	return &pesertas[0], nil
}

// Menyimpan atau mengupdate profil peserta milik akun user yang sedang login
func SaveMyPeserta(userID uint, input *model.Peserta) (*model.Peserta, error) {
	var existing model.Peserta
	err := config.DB.Where("user_id = ?", userID).First(&existing).Error

	tx := config.DB.Begin()

	if err != nil {
		// Belum pernah mengisi biodata peserta, buat baru
		input.UserID = &userID
		if err := tx.Create(input).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		// Sudah ada, perbarui nama
		if err := tx.Model(&existing).Update("nama", input.Nama).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// Hapus karakter formasi lama
		if err := tx.Where("peserta_id = ?", existing.ID).Delete(&model.PesertaKarakter{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// Masukkan karakter formasi baru
		for i := range input.Karakters {
			input.Karakters[i].PesertaID = existing.ID
		}
		if err := tx.Create(&input.Karakters).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		input.ID = existing.ID
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return GetPesertaByID(fmt.Sprintf("%d", input.ID))
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
