package entity

// User - Database Shape.
type User struct {
	ID       string
	Email    string
	Username string
	Image    string
	Bio      string
	Password string
	Role     Role

	Timestamps
}

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	return r == AdminRole || r == UserRole
}
