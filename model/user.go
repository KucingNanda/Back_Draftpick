package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id" example:"1"`
	Username string `gorm:"type:varchar(100);unique;not null" json:"username" example:"admin"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(50);not null;default:'member'" json:"role" example:"member"`
}

type AuthRequest struct {
	Username string `json:"username" example:""`
	Password string `json:"password" example:""`
}

type LoginRequest struct {
	Username string `json:"username" example:"member"`
	Password string `json:"password" example:"member123"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"admin123"`
	NewPassword string `json:"new_password" example:"admin456"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" example:"peserta"`
}
