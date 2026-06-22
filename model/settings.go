package model

type Settings struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	MaxCost     int    `gorm:"default:7" json:"maxCost"`
	Free4StarId   int    `gorm:"default:0" json:"free4StarId"`
	Free4StarHalf string `gorm:"default:'all'" json:"free4StarHalf"`
	BgmUrl        string `gorm:"default:''" json:"bgmUrl"`
}
