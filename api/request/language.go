package request

type LanguageRequest struct {
	Name string
}

type CreateLangRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Region      string `json:"region" validate:"required"`
	Description string `json:"description"`
}
