package service

import (
	"errors"
	"strings"

	"github.com/yugo412/leksikon/factory/dictionary"
	"github.com/yugo412/leksikon/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Dictionary struct {
	ORM *gorm.DB
	Log *zap.SugaredLogger

	Word   *model.WordModel
	Lang   *model.LanguageModel
	Source *model.SourceModel
}

func (d *Dictionary) Index(page int) (entries []model.Word, err error) {
	entries, err = d.Word.Index(page)

	return
}

func (d *Dictionary) Search(entry string) (word model.Word, err error) {
	entry = strings.TrimSpace(entry)

	word, err = d.Word.FindByEntry(entry)
	if errors.Is(gorm.ErrRecordNotFound, err) {
		// get registered 3rd party source before
		sources, _ := d.Source.GetSources()
		if len(sources) <= 0 {
			err = errors.New("no registered sources")

			return
		}

		// get entry from source
		// if it found on the first place, store it to the database and break process
		for _, s := range sources {
			dict, dictErr := dictionary.New(s.URL).Search(entry)
			if dictErr != nil {
				d.Log.Warnw("Can't find entry",
					"error", dictErr,
					"entry", entry,
				)
				continue
			} else {
				// store fetched entry to the db
				storeErr := d.Store(dict)
				if storeErr == nil {
					// re-fetch stored entry
					word, err = d.Word.FindByEntry(entry)
				}

				break
			}
		}
	}

	return
}

func (d *Dictionary) Store(dict dictionary.Word) (err error) {
	var word model.Word

	search := d.ORM.Where(model.Word{Entry: dict.Word}).First(&word)
	if errors.Is(search.Error, gorm.ErrRecordNotFound) {
		err = d.ORM.Transaction(func(ORM *gorm.DB) (err error) {
			lang, _ := d.Lang.DefaultLang()

			word = model.Word{
				LanguageID: uint64(lang.ID),
				Entry:      dict.Word,
				Spell:      dict.Spell,
				Syllable:   dict.Syllable,
				Source:     dict.Source,
			}

			result := ORM.Create(&word)
			if err != nil {
				return result.Error
			}

			// store slang words if exists
			if len(dict.Informals) > 0 {
				var slang []model.Slang
				for _, entry := range dict.Informals {
					slang = append(slang, model.Slang{
						WordID: uint64(word.ID),
						Entry:  entry,
					})
				}

				result = ORM.Create(&slang)
				if result.Error != nil {
					d.Log.Warnw("failed to store slang words",
						"error", result.Error,
						"entry", word.Entry,
					)
				}
			}

			for _, def := range dict.Definitions {
				definition := model.Definition{
					WordID:      uint64(word.ID),
					Description: def.Description,
				}

				result := ORM.Create(&definition)
				if result.Error != nil {
					d.Log.Errorw("failed to create definition",
						"error", result.Error,
						"word", word.Entry,
					)
				}

				if result.Error == nil {
					// define multi classes for each definition
					classes := map[string]model.Class{}
					if len(def.Classes) > 0 {
						classAlias := def.Classes[0]

						if class, ok := classes[classAlias]; !ok {
							result := d.ORM.Where(&model.Class{Alias: classAlias}).First(&class)
							if result.Error == nil {
								classes[classAlias] = class
							}
						}

						ORM.Create(&model.DefinitionClass{
							DefinitionID: uint64(definition.ID),
							ClassID:      uint64(classes[classAlias].ID),
						})
					}

					// create examples for each definition
					for _, example := range def.Examples {
						result := ORM.Create(&model.Example{
							DefinitionID: uint64(definition.ID),
							Sentence:     example,
						})

						if result.Error != nil {
							d.Log.Errorw("failed to create example",
								"error", result.Error,
								"word", word.Entry,
							)
						}
					}
				}
			}

			return
		})
	}

	return
}
