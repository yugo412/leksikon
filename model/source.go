package model

import "gorm.io/gorm"

type Source struct {
	Model

	Name     string `json:"name" gorm:"size:50"`
	URL      string `json:"url" gorm:"size:250"`
	IsActive bool   `json:"is_active"`
	Priority int    `json:"-"`
}

func (*Source) TableName() string {
	return "source"
}

type SourceModel struct {
	DB *gorm.DB
}

func (s *SourceModel) GetSources() (sources []*Source, err error) {
	err = s.DB.Model(&Source{}).
		Where("is_active = ?", true).
		Order("priority ASC").
		Find(&sources).
		Error

	return
}
