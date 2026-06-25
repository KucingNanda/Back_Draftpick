package repository

import (
	"draftpick-api/config"
	"draftpick-api/model"
)

func GetAllBgm() ([]model.Bgm, error) {
	var bgms []model.Bgm
	err := config.DB.Find(&bgms).Error
	return bgms, err
}

func GetBgmByID(id string) (model.Bgm, error) {
	var bgm model.Bgm
	err := config.DB.First(&bgm, id).Error
	return bgm, err
}

func CreateBgm(bgm *model.Bgm) error {
	return config.DB.Create(bgm).Error
}

func SaveBgm(bgm *model.Bgm) error {
	return config.DB.Save(bgm).Error
}

func DeleteBgm(id string) error {
	return config.DB.Delete(&model.Bgm{}, id).Error
}
