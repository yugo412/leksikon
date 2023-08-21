package model

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Example struct {
	Model

	DefinitionID uint64 `json:"-"`
	Sentence     string `json:"sentence" gorm:"type:text"`
}

func (*Example) TableName() string {
	return "example"
}

func (e *Example) BeforeCreate(_ *gorm.DB) error {
	e.ULID = ulid.Make().String()

	return nil
}

type ExampleModel struct {
	ORM *gorm.DB
}
