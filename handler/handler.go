package handler

import (
	"fmt"

	"github.com/FilmListClub/backend/dao"
)

type Handler struct {
	dao dao.DAO
}

func New(dao dao.DAO) *Handler {
	return &Handler{dao}
}

// --- Errors ---

type Error struct {
	Message      string
	ResponseCode int
}

// Wrap mutates an existing error message to add context and returns it
func (e *Error) Wrap(message string) *Error {
	e.Message = fmt.Sprintf("%s: %s", message, e.Message)
	return e
}
