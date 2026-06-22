package model

import "time"

type Peserta struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	Nama      string            `gorm:"type:varchar(100);not null" json:"nama"`
	Karakters []PesertaKarakter `gorm:"foreignKey:PesertaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"karakters"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type PesertaKarakter struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	PesertaID  uint     `gorm:"not null;uniqueIndex:idx_peserta_karakter" json:"pesertaId"`
	KarakterID uint     `gorm:"not null;uniqueIndex:idx_peserta_karakter" json:"karakterId"`
	Karakter   Karakter `gorm:"foreignKey:KarakterID" json:"karakter"`
	IsC1       bool     `gorm:"default:false" json:"isC1"`
}
