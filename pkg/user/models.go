package user

import "time"

// User is a registered user.
type User struct {
	ID        string
	Phone     string
	Secret    string
	FirstName string
	LastName  string
	Email     string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
