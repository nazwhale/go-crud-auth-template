package dao

import (
	"database/sql"
)

type DAO interface {
	SaveUser(User) error
	ReadUserByEmail(string) (User, error)

	SaveListItem(ListItem) error
	ReadListItemsForUser(int) ([]ListItem, error)
}

type dao struct {
	db *sql.DB
}

func New(db *sql.DB) DAO {
	return &dao{db}
}
