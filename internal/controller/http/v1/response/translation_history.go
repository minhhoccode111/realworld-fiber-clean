package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// TranslationHistory -.
// Try to extend the TranslationHistory example with pagination and put inside
// response package instead of entity
type TranslationHistory struct {
	History []entity.TranslationClone `json:"history"`

	entity.Pagination
}
