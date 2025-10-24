package entity

// User -.
type User struct {
	Id       string
	Email    string
	Username string
	Image    string
	Bio      string
	Password string

	Timestamps
}
