package dao

import (
	"errors"
	"fmt"
	"time"
)

type ListItem struct {
	ID        int
	UserID    int
	Title     string
	CreatedAt time.Time
}

func (d dao) SaveListItem(item ListItem) error {
	t := time.Now().UTC()

	sqlStatement := `
INSERT INTO listitems (user_id, title, created_at)
VALUES ($1, $2, $3)
RETURNING id`

	var id int
	if err := d.db.QueryRow(
		sqlStatement,
		item.UserID,
		item.Title,
		t,
	).Scan(&id); err != nil {
		return errors.New(fmt.Sprintf("error writing list item to db: %q", err))
	}

	return nil
}

func (d dao) ReadListItemsForUser(userID int) ([]ListItem, error) {
	sqlStatement := `
SELECT * FROM listitems
WHERE user_id=$1
`
	var items []ListItem
	rows, err := d.db.Query(sqlStatement, userID)
	if err != nil {
		return []ListItem{}, errors.New("error querying db")
	}
	defer rows.Close()

	for rows.Next() {
		var item ListItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.Title, &item.CreatedAt); err != nil {
			return []ListItem{}, errors.New("error scanning row")
		}

		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return []ListItem{}, errors.New("error scanning row")
	}

	return items, nil
}
