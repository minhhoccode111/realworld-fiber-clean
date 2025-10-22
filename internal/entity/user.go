package entity

// User -.
type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
	Password string `json:"-"`

	Timestamps
}
