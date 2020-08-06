package dao

import (
	"database/sql"
)

// See https://github.com/nazwhale/go-crud-starter/blob/master/src/main.go

type DAO interface {
	SaveUser(User) error
	ReadUserByEmail(string) (User, error)

	SaveListItem(ListItem) error
	ReadListItemsForUser(int) ([]ListItem, error)
}

type dao struct {
	db *sql.DB
}

// Returns the db wrapped in a dao object
func New(db *sql.DB) DAO {
	return &dao{db}
}
