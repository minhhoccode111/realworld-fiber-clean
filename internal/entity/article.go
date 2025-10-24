package entity

// Article -.
type Article struct {
	Id          string
	AuthorId    string
	Slug        string
	Title       string
	Body        string
	Description string

	Timestamps
}
