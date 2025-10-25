package entity

// Article -. Database Shape
type Article struct {
	Id          string
	AuthorId    string
	Slug        string
	Title       string
	Body        string
	Description string

	Timestamps
}

type ArticleDetail struct {
	Slug           string         `json:"slug"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	TagList        []string       `json:"tagList"`
	Favorited      bool           `json:"favorited"`
	FavoritesCount int            `json:"favoritesCount"`
	Author         ProfilePreview `json:"author"`
	Body           string         `json:"body"`

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
	// Body           string         `json:"body"`

	Timestamps
}
