package model

import (
	"gorm.io/gorm"

	"github.com/oklog/ulid/v2"
)

type Class struct {
	Model

	Alias string `json:"alias" gorm:"size:20"`
	Name  string `json:"name" gorm:"size:50"`

	Definitions []Definition `json:"definitions,omitempty" gorm:"many2many:definition_class;"`
}

func (*Class) TableName() string {
	return "class"
}

func (c *Class) BeforeCreate(_ *gorm.DB) error {
	c.ULID = ulid.Make().String()

	return nil
}
