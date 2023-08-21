package model

import (
	"errors"
	"gorm.io/gorm"
)

const (
	DefaultLanguage = "Indonesia"
	DefaultRegion   = "Indonesia"
)

type Language struct {
	Model

	Alias  string `json:"alias" gorm:"size:20;index"`
	Name   string `json:"name" gorm:"index;type:varchar(100);not null"`
	Region string `json:"-" gorm:"size:250"`

	Words []Word `json:"-"`
}

func (*Language) TableName() string {
	return "language"
}

type LanguageModel struct {
	ORM *gorm.DB
}

func (l *LanguageModel) DefaultLang() (lang Language, err error) {
	err = l.ORM.Where(&Language{Name: DefaultLanguage}).First(&lang).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		lang = Language{
			Name:   DefaultLanguage,
			Region: DefaultRegion,
		}
		err = l.ORM.Create(&lang).Error
	}

	return
}
