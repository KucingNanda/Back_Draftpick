package model

type Bgm struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Nama string `json:"nama"`
	Url  string `json:"url"`
}
