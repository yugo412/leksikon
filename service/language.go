package service

import (
	"github.com/yugo412/leksikon/api/request"
	"github.com/yugo412/leksikon/library"
	"github.com/yugo412/leksikon/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Language struct {
	ORM *gorm.DB
	Log *zap.SugaredLogger

	Lang *model.LanguageModel
}

func (l *Language) FetchLocalLanguage() (err error) {
	locals, err := library.FetchLocalLanguage()
	if err != nil {
		return
	}

	var languages []model.Language
	for _, local := range locals {
		languages = append(languages, model.Language{
			Name:   local.Name,
			Region: local.Region,
		})
	}

	err = l.ORM.Create(&languages).Error

	return
}

func (l *Language) GetAllLanguage(req request.LanguageRequest) (languages []model.Language, err error) {
	err = l.ORM.
		Where("name LIKE ?", "%"+req.Name+"%").
		Find(&languages).
		Error

	return
}
