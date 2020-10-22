package client

import (
	"database/sql"

	"github.com/aslrousta/rand"
	"github.com/aslrousta/sso/lib/service"
)

const (
	clientIDLen     = 8
	clientSecretLen = 16
)

// Register creates a new client.
func Register(db *sql.DB, ownerID, name, domain string) (*Client, error) {
	clientID, err := rand.RandomString(clientIDLen, rand.Digit)
	if err != nil {
		return nil, &service.Error{
			Service: "client.Register",
			Message: "faild to generate client id",
			Cause:   err,
		}
	}

	clientSecret, err := rand.RandomString(clientSecretLen, rand.All)
	if err != nil {
		return nil, &service.Error{
			Service: "client.Register",
			Message: "faild to generate client secret",
			Cause:   err,
		}
	}

	c := Client{
		ID:      clientID,
		Secret:  clientSecret,
		OwnerID: ownerID,
		Name:    name,
		Domain:  domain,
	}
	err = db.QueryRow(
		"INSERT INTO clients (id,secret,owner_id,name,domain)"+
			" VALUES ($1,$2,$3,$4,$5)"+
			" RETURNING created_at, updated_at",
		c.ID,
		c.Secret,
		c.OwnerID,
		c.Name,
		c.Domain,
	).Scan(
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, &service.Error{
			Service: "client.Register",
			Message: "failed to store client",
			Cause:   err,
		}
	}
	return &c, nil
}

// FindByID retrives a client by its ID.
func FindByID(db *sql.DB, id string) (*Client, error) {
	var c Client
	err := db.QueryRow(
		"SELECT id,secret,owner_id,name,domain,created_at,updated_at"+
			" FROM clients"+
			" WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Scan(
		&c.ID,
		&c.Secret,
		&c.OwnerID,
		&c.Name,
		&c.Domain,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &service.NotFoundError{Key: id}
		}
		return nil, &service.Error{
			Service: "client.FindByID",
			Message: "failed to retrieve the client",
			Cause:   err,
		}
	}
	return &c, nil
}

// Update updates a client's information.
func Update(db *sql.DB, id, name, domain string) error {
	_, err := db.Exec(
		"UPDATE clients SET name = $1, domain = $2 WHERE id = $3",
		name,
		domain,
		id,
	)
	if err != nil {
		return &service.Error{
			Service: "client.Update",
			Message: "failed to update the client",
			Cause:   err,
		}
	}
	return nil
}

// Delete removes a client.
func Delete(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE clients SET deleted_at = NOW() WHERE id = $1", id)
	if err != nil {
		return &service.Error{
			Service: "client.Delete",
			Message: "failed to delete the client",
			Cause:   err,
		}
	}
	return nil
}
