package model

type DefinitionClass struct {
	DefinitionID uint64 `json:"-"`
	ClassID      uint64 `json:"-"`
}

func (*DefinitionClass) TableName() string {
	return "definition_class"
}
