package web

import (
	"html/template"
	"net/http"

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

	h.Use(middleware.Logger)
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.ThreadsList())
	})

	return h
}

const threadsListHTML = `
<h1>Threads</h1>
<dl>
{{range .Threads}}
	<dt><strong>{{.Title}}</strong></dt>
	<dt>{{.Description}}</dt>
{{end}}
</dl>
`

func (h *Handler) ThreadsList() http.HandlerFunc {
	type data struct {
		Threads []miji.Thread
	}

	tmpl := template.Must(template.New("").Parse(``))

	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Threads: tt})
	}
}
