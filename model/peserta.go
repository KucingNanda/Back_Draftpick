package model

import "time"

type Peserta struct {
	ID        uint              `gorm:"primaryKey" json:"id" example:"1"`
	UserID    *uint             `gorm:"uniqueIndex" json:"userId,omitempty" example:"2"`
	User      *User             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Nama      string            `gorm:"type:varchar(100);not null" json:"nama" example:"Nanda"`
	Karakters []PesertaKarakter `gorm:"foreignKey:PesertaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"karakters"`
	CreatedAt time.Time         `json:"createdAt" example:"2026-06-26T12:00:00Z"`
	UpdatedAt time.Time         `json:"updatedAt" example:"2026-06-26T12:00:00Z"`
}

type PesertaKarakter struct {
	ID         uint     `gorm:"primaryKey" json:"id" example:"1"`
	PesertaID  uint     `gorm:"not null;uniqueIndex:idx_peserta_karakter" json:"pesertaId" example:"1"`
	KarakterID uint     `gorm:"not null;uniqueIndex:idx_peserta_karakter" json:"karakterId" example:"5"`
	Karakter   Karakter `gorm:"foreignKey:KarakterID" json:"karakter"`
	IsC1       bool     `gorm:"default:false" json:"isC1" example:"false"`
}
