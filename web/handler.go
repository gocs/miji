package web

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gocs/miji"
)

type Handler struct {
	*chi.Mux

	store miji.Store
}

func NewHandler(store miji.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	threads := ThreadHandler{store: store}
	posts := PostHandler{store: store}
	// comments := CommentHandler{store: store}
	// users:= UserHandler{store: store}

	h.Use(middleware.Logger)
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", threads.List())
		r.Get("/new", threads.Create())
		r.Post("/", threads.Store())
		r.Get("/{id}", threads.Show())
		r.Post("/{id}/delete", threads.Delete())
		r.Post("/update", threads.Update())
		r.Get("/{id}/new", posts.Create())
		r.Post("/{id}", posts.Store())
		r.Get("/{threadID}/{postID}", posts.Show())
		r.Post("/{threadID}/{postID}/delete", posts.Delete())
		r.Post("/{threadID}/update", posts.Update())
		r.Post("/{threadID}/{postID}/vote", posts.Vote())
		// r.Post("/{threadID}/{postID}", comments.Store())
	})
	// h.Get("/comments/{id}/vote", comments.Vote())
	// h.Get("/register", users.Register())
	// h.Post("/register", users.RegisterSubmit())
	// h.Get("/login", users.Login())
	// h.Post("/login", users.LoginSubmit())
	// h.Get("/logout", users.Logout())

	return h
}
