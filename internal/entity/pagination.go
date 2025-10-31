package entity

type Pagination struct {
	Limit  uint64 `json:"limit"         example:"10"`
	Offset uint64 `json:"offset"        example:"0"`
	// should be "total" but RealWorld API Specs expect "articlesCount"
	Total uint64 `json:"articlesCount" example:"18"`
}
