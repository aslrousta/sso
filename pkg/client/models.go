package client

import "time"

// Client is a registered client application.
type Client struct {
	ID        string
	Secret    string
	OwnerID   string
	Name      string
	Domain    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
