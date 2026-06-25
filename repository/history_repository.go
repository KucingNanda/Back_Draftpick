package repository

import (
	"draftpick-api/config"
	"draftpick-api/model"
)

func SaveHistory(history *model.History) error {
	return config.DB.Create(history).Error
}

func SaveExistingHistory(history *model.History) error {
	return config.DB.Save(history).Error
}

func GetHistoryCountByRole(role string) (int64, error) {
	var count int64
	err := config.DB.Model(&model.History{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

func GetOldestHistoryByRole(role string) (model.History, error) {
	var oldest model.History
	err := config.DB.Where("role = ?", role).Order("updated_at asc").First(&oldest).Error
	return oldest, err
}

func GetHistoriesByRole(role string) ([]model.History, error) {
	var histories []model.History
	err := config.DB.Where("role = ?", role).Order("updated_at desc").Find(&histories).Error
	return histories, err
}

func DeleteHistory(id string) error {
	return config.DB.Delete(&model.History{}, id).Error
}
