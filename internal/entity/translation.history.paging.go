package entity

// TranslationHistoryPaging -.
// Try to extend the TranslationHistory example with pagination
type TranslationHistoryPaging struct {
	History []Translation `json:"history"`
	Limit   int           `json:"limit"`
	Offset  int           `json:"offset"`
}
