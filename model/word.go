package model

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Word struct {
	Model

	LanguageID uint64 `json:"-"`
	Entry      string `json:"word" gorm:"index;type:varchar(150);not null"`
	Spell      string `json:"spell" gorm:"size:150"`
	Syllable   string `json:"syllable" gorm:"size:50"`
	Source     string `json:"-" gorm:"size:150"`

	// custom fields to simplify response
	Informals []string `json:"informals" gorm:"-"`
	Counts    int64    `json:"views,omitempty" gorm:"-"`

	Language    Language     `json:"-"`
	Slang       []Slang      `json:"-"`
	Definitions []Definition `json:"definitions,omitempty"`
	Counters    []Counter    `json:"-"`
}

func (*Word) TableName() string {
	return "word"
}

type WordModel struct {
	ORM *gorm.DB
}

func (w *Word) BeforeCreate(_ *gorm.DB) error {
	w.ULID = ulid.Make().String()

	return nil
}

func (w *WordModel) Index(page int) (entries []Word, err error) {
	limit := 10
	offset := (page - 1) * limit

	err = w.ORM.
		Select("id, entry, syllable, spell").
		Preload("Slang").
		Order("entry ASC").
		Offset(offset).
		Limit(limit).
		Find(&entries).
		Error

	return
}

func (w *WordModel) FindByEntry(entry string) (word Word, err error) {
	err = w.ORM.
		Preload("Slang").
		Preload("Language").
		Preload("Definitions", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Examples").
				Preload("Classes")
		}).
		Joins(fmt.Sprintf("LEFT JOIN %s AS i ON i.word_id = word.id", (*Slang)(nil).TableName())).
		Where("word.entry = ? OR i.entry = ?", entry, entry).
		First(&word).
		Error

	word.Informals = make([]string, 0)
	for _, s := range word.Slang {
		word.Informals = append(word.Informals, s.Entry)
	}

	for i := range word.Definitions {
		d := &word.Definitions[i]
		d.Sentences = make([]string, 0)

		for _, e := range d.Examples {
			d.Sentences = append(d.Sentences, e.Sentence)
		}
	}

	return

}
