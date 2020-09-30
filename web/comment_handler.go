package web

import (
	"net/http"

	"github.com/gocs/miji"
)

type CommentHandler struct {
	store miji.Store
}

func (h *CommentHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
