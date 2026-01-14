package entity

// Comment represents the stored database shape for a comment.
type Comment struct {
	ID        string
	ArticleID string
	AuthorID  string
	Body      string

	Timestamps
}

// CommentDetail contains comment information returned to clients.
type CommentDetail struct {
	ID     string         `json:"id"`
	Body   string         `json:"body"`
	Author ProfilePreview `json:"author"`

	Timestamps
}
