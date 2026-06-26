package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id" example:"1"`
	Username string `gorm:"type:varchar(100);unique;not null" json:"username" example:"admin"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(50);not null;default:'peserta'" json:"role" example:"wasit"`
}

type AuthRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
	Role     string `json:"role" example:"wasit"`
}

type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"admin123"`
	NewPassword string `json:"new_password" example:"admin456"`
}
