package response

type WordDefinition struct {
	Description string   `json:"description"`
	Classes     []string `json:"classes"`
	Examples    []string `json:"examples"`
}

type WordResponse struct {
	Root        string           `json:"root"`
	Word        string           `json:"word"`
	Spell       string           `json:"spell"`
	Syllable    string           `json:"syllable"`
	Informal    string           `json:"informal"`
	Definitions []WordDefinition `json:"definitions"`
}
