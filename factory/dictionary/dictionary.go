package dictionary

import (
	"errors"
)

type Dictionary interface {
	Vendor() string
	Search(string) (Word, error)
}

type Word struct {
	Root        string
	Word        string
	Spell       string
	Syllable    string
	Source      string
	Informals   []string
	Definitions []WordDefinition
}

type WordDefinition struct {
	Description string
	Classes     []string
	Examples    []string
}

var (
	ErrMaxLimit = errors.New("max limit reached")
)

func New(source string) Dictionary {
	return &KBBI{
		URL: source,
	}
}
