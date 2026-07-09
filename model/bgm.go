package model

type Bgm struct {
	ID   uint   `gorm:"primaryKey" json:"id" example:"1"`
	Nama string `gorm:"type:varchar(255);not null" json:"nama" example:"Let the Wind Tell You"`
	Url  string `gorm:"type:text;not null" json:"url" example:"https://res.cloudinary.com/.../audio.mp3"`
}
