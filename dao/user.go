package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrNoUserExists = errors.New("no user exists with email")

type User struct {
	Email          string
	HashedPassword string
}

func (d dao) SaveUser(user User) error {
	t := time.Now().UTC()
	sqlStatement := `
INSERT INTO users (email, password, created_at)
VALUES ($1, $2, $3)
RETURNING id`

	var id int
	if err := d.db.QueryRow(
		sqlStatement,
		user.Email,
		user.HashedPassword,
		t,
	).Scan(&id); err != nil {
		return errors.New(fmt.Sprintf("error writing user to db: %q", err))
	}

	return nil
}

func (d dao) ReadUserByEmail(email string) (User, error) {
	sqlStatement := `
SELECT email, password from users
WHERE email=$1`

	var user User
	result := d.db.QueryRow(sqlStatement, email)
	err := result.Scan(&user.Email, &user.HashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrNoUserExists
		}
		return User{}, errors.New(fmt.Sprintf("error reading user by email: %q", err))
	}

	return user, nil
}
