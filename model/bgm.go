package model

type Bgm struct {
	ID   uint   `gorm:"primaryKey" json:"id" example:"1"`
	Nama string `json:"nama" example:"Let the Wind Tell You"`
	Url  string `json:"url" example:"https://res.cloudinary.com/.../audio.mp3"`
}
