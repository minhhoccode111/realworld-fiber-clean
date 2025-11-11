package entity

type Pagination struct {
	Limit  uint64 `json:"limit"         example:"10"`
	Offset uint64 `json:"offset"        example:"0"`
	Total  uint64 `json:"total"         example:"18"`
	Dummy  uint64 `json:"articlesCount"`
	// 'Dummy' is just dummy because realworld api specs need its presence
}
