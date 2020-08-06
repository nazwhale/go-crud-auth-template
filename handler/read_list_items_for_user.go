package handler

import (
	"encoding/json"
	"net/http"
)

type readListItemsForUserReq struct {
	UserID int `json:"user_id"`
}

func (h *Handler) ReadListItemsForUser(w http.ResponseWriter, r *http.Request) *Error {
	req := &readListItemsForUserReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusBadRequest,
		}
	}

	if req.UserID == 0 {
		return &Error{
			Message:      "userID is empty",
			ResponseCode: http.StatusBadRequest,
		}
	}

	items, err := h.dao.ReadListItemsForUser(req.UserID)
	if err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
		return e.Wrap("error reading list items from db")
	}

	json.NewEncoder(w).Encode(items)

	return nil
}
