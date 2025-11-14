package entity

// User - Database Shape.
type User struct {
	ID       string
	Email    string
	Username string
	Image    string
	Bio      string
	Password string

	Timestamps
}
