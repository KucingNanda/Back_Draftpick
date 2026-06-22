package model

// Karakter merepresentasikan tabel karakters di PostgreSQL
type Karakter struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nama     string `gorm:"type:varchar(100);not null" json:"nama"`
	Vision   string `gorm:"type:varchar(50)" json:"vision"`   // Anemo, Geo, Electro, Dendro, Hydro, Pyro, Cryo, Flex
	Senjata  string `gorm:"type:varchar(50)" json:"senjata"`  // Sword, Claymore, Catalyst, Bow, Polearm
	Rarity   string `gorm:"type:varchar(50)" json:"rarity"`   // 4 Star, 5 Star Limited, dll
	Cost     int    `gorm:"type:int" json:"cost"`             // 0 atau 1
	ImageURL string `gorm:"type:text" json:"imageUrl"`        // Link Cloudinary
}
