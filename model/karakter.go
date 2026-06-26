package model

// Karakter merepresentasikan tabel karakters di PostgreSQL
type Karakter struct {
	ID         uint   `gorm:"primaryKey" json:"id" example:"1"`
	Nama       string `gorm:"type:varchar(100);not null" json:"nama" example:"Zhongli"`
	Elemen     string `gorm:"type:varchar(50)" json:"elemen" example:"Geo"`
	Senjata    string `gorm:"type:varchar(50)" json:"senjata" example:"Polearm"`
	Rarity     string `gorm:"type:varchar(50)" json:"rarity" example:"5 Star Limited"`
	Cost       int    `gorm:"type:int" json:"cost" example:"1"`
	ImageURL   string `gorm:"type:text" json:"imageUrl" example:"https://res.cloudinary.com/.../zhongli.png"`
	Konstelasi string `gorm:"type:varchar(100)" json:"konstelasi" example:"Lapis Dei"`
	UlangTahun string `gorm:"type:varchar(50)" json:"ulangTahun" example:"31 Desember"`
	Region     string `gorm:"type:varchar(50);default:'Tidak Ada'" json:"region" example:"Liyue"`
	Sinergi    string `gorm:"type:varchar(50);default:'Tidak Ada'" json:"sinergi" example:"Tidak Ada"`
}
