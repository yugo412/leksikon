package model

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Definition struct {
	Model

	WordID      uint64 `json:"-"`
	Description string `json:"description" gorm:"type:text"`

	// custom fields to simplify response
	Sentences []string `json:"examples" gorm:"-"`

	Word     Word      `json:"-"`
	Classes  []Class   `json:"classes" gorm:"many2many:definition_class;"`
	Examples []Example `json:"-"`
}

func (*Definition) TableName() string {
	return "definition"
}

func (d *Definition) BeforeCreate(_ *gorm.DB) error {
	d.ULID = ulid.Make().String()

	return nil
}

type DefinitionModel struct {
	ORM *gorm.DB
}
