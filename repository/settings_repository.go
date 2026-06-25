package repository

import (
	"draftpick-api/config"
	"draftpick-api/model"
)

func GetSettings() (model.Settings, error) {
	var settings model.Settings
	err := config.DB.First(&settings).Error
	return settings, err
}

func SaveSettings(settings *model.Settings) error {
	return config.DB.Save(settings).Error
}
