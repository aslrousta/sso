package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/aslrousta/rand"
	"github.com/aslrousta/sso/lib/service"
	"github.com/aslrousta/sso/lib/types"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

const (
	idLen    = 16
	tokenLen = 8
	codeLen  = 5
	codeTTL  = 5 * time.Minute
	columns  = "id,phone,secret,first_name,last_name,email,bio,created_at,updated_at"
)

var (
	errUserNotFound = errors.New("user not found")
)

// Register initiates authentication for a phone number.
func Register(db *sql.DB, r *redis.Client, phone types.Cellphone) (token, code string, err error) {
	token, err = rand.RandomString(tokenLen, rand.All)
	if err != nil {
		return "", "", &service.Error{
			Service: "user.Register",
			Message: "failed to generate token",
			Cause:   err,
		}
	}
	code, err = rand.RandomString(codeLen, rand.Digit)
	if err != nil {
		return "", "", &service.Error{
			Service: "user.Register",
			Message: "failed to generate auth code",
			Cause:   err,
		}
	}
	values := url.Values{}
	values.Set("phone", string(phone))
	values.Set("code", code)
	err = r.Set(
		context.Background(),
		fmt.Sprintf("reg:%s", token), values.Encode(),
		codeTTL,
	).Err()
	if err != nil {
		return "", "", &service.Error{
			Service: "user.Register",
			Message: "failed to store auth info",
			Cause:   err,
		}
	}
	return token, code, nil
}

// Authenticate signs a new user up or signs an existing user in.
func Authenticate(db *sql.DB, r *redis.Client, token, code string) (*User, error) {
	res, err := r.Get(context.Background(), fmt.Sprintf("reg:%s", token)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, &service.NotFoundError{
				Entity: "auth info",
				Key:    token,
			}
		}
		return nil, &service.Error{
			Service: "user.Authenticate",
			Message: "failed to find auth info",
			Cause:   err,
		}
	}
	values, err := url.ParseQuery(res)
	if err != nil {
		return nil, &service.Error{
			Service: "user.Authenticate",
			Message: "failed to decode auth info",
			Cause:   err,
		}
	}
	if values.Get("code") != code {
		return nil, &service.NotFoundError{
			Entity: "auth info",
			Key:    token,
		}
	}
	phone := types.Cellphone(values.Get("phone"))
	u, err := find(db, phone)
	if err != nil {
		if err != errUserNotFound {
			return nil, &service.Error{
				Service: "user.Authenticate",
				Message: "faild to find user",
				Cause:   err,
			}
		}
		u, err = create(db, phone)
		if err != nil {
			return nil, &service.Error{
				Service: "user.Authenticate",
				Message: "failed to create user",
				Cause:   err,
			}
		}
	}
	return u, nil
}

func create(db *sql.DB, phone types.Cellphone) (*User, error) {
	userID, err := rand.RandomString(8, rand.Digit)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(phone), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := User{
		ID:     userID,
		Phone:  phone.Masked(),
		Secret: string(hash),
	}
	err = db.QueryRow(
		"INSERT INTO users (id,phone,secret,created_at,updated_at)"+
			" VALUES ($1,$2,$3,NOW(),NOW())"+
			" RETURNING created_at,updated_at",
		u.ID, u.Phone, u.Secret,
	).Scan(
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func find(db *sql.DB, phone types.Cellphone) (*User, error) {
	rows, err := db.Query(
		"SELECT "+columns+" FROM users WHERE phone = $1",
		phone.Masked(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := User{}
	found := false
	for rows.Next() {
		err := rows.Scan(
			&u.ID, &u.Phone, &u.Secret, &u.FirstName, &u.LastName, &u.Email,
			&u.Bio, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		err = bcrypt.CompareHashAndPassword([]byte(u.Secret), []byte(phone))
		if err == nil {
			found = true
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !found {
		return nil, errUserNotFound
	}
	return &u, nil
}
