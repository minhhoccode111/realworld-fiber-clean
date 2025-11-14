package entity

// Article -. Database Shape.
type Article struct {
	ID          string
	AuthorID    string
	Slug        string
	Title       string
	Body        string
	Description string

	Timestamps
}

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

type ArticleDetail struct {
	Body string `json:"body"`

	ArticlePreview
	Timestamps
}
