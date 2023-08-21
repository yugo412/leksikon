package model

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Slang struct {
	Model

	WordID uint64 `json:"-"`
	Entry  string `json:"entry" gorm:"index;size:50"`
}

func (*Slang) TableName() string {
	return "informal"
}

func (s *Slang) BeforeCreate(_ *gorm.DB) error {
	s.ULID = ulid.Make().String()

	return nil
}

type SlangModel struct {
	ORM *gorm.DB
}
