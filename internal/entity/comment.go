package entity

// Comment -. Database Shape
type Comment struct {
	Id        string
	ArticleId string
	AuthorId  string
	Body      string

	Timestamps
}

type CommentDetail struct {
	Id     string         `json:"id"`
	Body   string         `json:"body"`
	Author ProfilePreview `json:"author"`

	Timestamps
}
