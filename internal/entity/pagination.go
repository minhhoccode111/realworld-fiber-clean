package entity

type Pagination struct {
	Limit  uint64 `json:"limit"         example:"10"`
	Offset uint64 `json:"offset"        example:"0"`
	Total  uint64 `json:"total"         example:"18"`
	Dummy  uint64 `json:"articlesCount"`
	/*
		Realworld API test script need 'articlesCount' property, but I want
		Pagination to be reusable accross Articles, Comments, Tags etc. So the
		real one is 'total' and 'articlesCount' is just dummy field
	*/
}
