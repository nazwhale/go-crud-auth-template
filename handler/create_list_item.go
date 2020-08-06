package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FilmListClub/backend/dao"
)

type createListItemReq struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
}

func (h *Handler) CreateListItem(w http.ResponseWriter, r *http.Request) *Error {
	req := &createListItemReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusBadRequest,
		}
	}

	switch {
	case req.UserID == 0:
		return &Error{
			Message:      "userID is empty",
			ResponseCode: http.StatusBadRequest,
		}
	case req.Title == "":
		return &Error{
			Message:      "title is empty",
			ResponseCode: http.StatusBadRequest,
		}
	}

	listItem := dao.ListItem{
		UserID: req.UserID,
		Title:  req.Title,
	}

	if err := h.dao.SaveListItem(listItem); err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}

		return e.Wrap("error saving list item to db")
	}

	return nil
}
