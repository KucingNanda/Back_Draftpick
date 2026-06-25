package repository

import (
	"draftpick-api/config"
	"draftpick-api/model"
)

func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func GetUserByID(id interface{}) (model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return user, err
}

func SaveUser(user *model.User) error {
	return config.DB.Save(user).Error
}

func GetUsersByRole(role string) ([]model.User, error) {
	var users []model.User
	err := config.DB.Select("id", "username", "role").Where("role = ?", role).Find(&users).Error
	return users, err
}
