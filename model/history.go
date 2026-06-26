package model

import "time"

type History struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Role      string    `json:"role" example:"wasit"` // "wasit" atau "peserta"
	PlayerA   string    `json:"playerA" example:"Nanda"`
	PlayerB   string    `json:"playerB" example:"Rizky"`
	DraftData string    `json:"draftData" gorm:"type:text" example:"{\"turn\": 1, \"phase\": \"ban\"}"` // Berisi JSON DraftState
	CreatedAt time.Time `json:"created_at" example:"2026-06-26T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-06-26T12:00:00Z"`
}
