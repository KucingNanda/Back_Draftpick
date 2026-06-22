package model

import "time"

type History struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Role      string    `json:"role"` // "wasit" atau "peserta"
	PlayerA   string    `json:"playerA"`
	PlayerB   string    `json:"playerB"`
	DraftData string    `json:"draftData" gorm:"type:text"` // Berisi JSON DraftState
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
