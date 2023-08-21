package model

import "time"

type Counter struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	WordID    uint64    `json:"-"`
	CreatedAt time.Time `json:"-"`

	Word Word `json:"words"`
}

func (*Counter) TableName() string {
	return "counter"
}
