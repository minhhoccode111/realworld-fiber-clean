package entity

// TranslationHistoryPaging -.
// Try to extend the TranslationHistory example with pagination
type TranslationCloneHistory struct {
	History []TranslationClone `json:"history"`
	Limit   uint64             `json:"limit"   example:"10"`
	Offset  uint64             `json:"offset"  example:"0"`
	Total   uint64             `json:"total"   example:"20"`
}
