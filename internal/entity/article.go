package entity

// Article represents the stored database shape for an article.
type Article struct {
	ID          string
	AuthorID    string
	Slug        string
	Title       string
	Body        string
	Description string

	Timestamps
}

// ArticlePreview contains article information for list responses.
type ArticlePreview struct {
	Slug           string         `json:"slug"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	TagList        []string       `json:"tagList"`
	Favorited      bool           `json:"favorited"`
	FavoritesCount int            `json:"favoritesCount"`
	Author         ProfilePreview `json:"author"`

	Timestamps
}

// ArticleDetail extends the preview with the full body content.
type ArticleDetail struct {
	Body string `json:"body"`

	ArticlePreview
	Timestamps
}
