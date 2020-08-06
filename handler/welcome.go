package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) *Error {
	w.Write([]byte(fmt.Sprintf("👋 Welcome: %s", "[name]")))
	return nil
}
