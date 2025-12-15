package entity

// User represents a stored user entity.
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

// Role defines the authorization role assigned to a user.
type Role string

const (
	// AdminRole grants administrative privileges.
	AdminRole Role = "admin"
	// UserRole grants standard user access.
	UserRole Role = "user"
)

// String returns the role value as a string.
func (r Role) String() string {
	return string(r)
}

// IsValid reports whether the role matches one of the allowed values.
func (r Role) IsValid() bool {
	return r == AdminRole || r == UserRole
}
