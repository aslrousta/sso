package client

import "time"

// Client is a registered client application.
type Client struct {
	ID        string
	Name      string
	Secret    string
	Domain    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
