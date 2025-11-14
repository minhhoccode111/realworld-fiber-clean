package entity

// Comment -. Database Shape.
type Comment struct {
	ID        string
	ArticleID string
	AuthorID  string
	Body      string

	Timestamps
}

type CommentDetail struct {
	ID     string         `json:"id"`
	Body   string         `json:"body"`
	Author ProfilePreview `json:"author"`

	Timestamps
}
