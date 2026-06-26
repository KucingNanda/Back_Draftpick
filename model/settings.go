package model

type Settings struct {
	ID          uint   `gorm:"primaryKey" json:"id" example:"1"`
	MaxCost     int    `gorm:"default:7" json:"maxCost" example:"7"`
	Free4StarId   int    `gorm:"default:0" json:"free4StarId" example:"0"`
	Free4StarHalf string `gorm:"default:'all'" json:"free4StarHalf" example:"all"`
	BgmUrl        string `gorm:"default:''" json:"bgmUrl" example:"https://res.cloudinary.com/.../audio.mp3"`
}
