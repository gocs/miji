package web

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gocs/miji"
	"github.com/google/uuid"
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
		r.Get("/new", h.ThreadsCreate())
		r.Post("/", h.ThreadsStore())
		r.Post("/{id}/delete", h.ThreadsDelete())
		r.Post("/{id}/update", h.ThreadsUpdatePage())
		r.Post("/update", h.ThreadsUpdate())
	})

	return h
}

const threadsListHTML = `
<h1>Threads</h1>
<dl>{{range .Threads}}
	<dt><strong>{{.Title}}</strong></dt>
	<dd>{{.Description}}</dd>
	<dd><form action="/threads/{{.ID}}/update" method="POST">
			<button type="submit">Update</button>
		</form></dd>
	<dd><form action="/threads/{{.ID}}/delete" method="POST">
			<button type="submit">Delete</button>
		</form></dd>
{{end}}</dl>
<a href="/threads/new"> Create New</a>
`

func (h *Handler) ThreadsList() http.HandlerFunc {
	type data struct {
		Threads []miji.Thread
	}

	tmpl := template.Must(template.New("").Parse(threadsListHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Threads: tt})
	}
}

const threadsCreateHTML = `
<h1>New Thread</h1>
<form action="/threads" method="POST">
	<table>
		<tr>
			<td>Title</td>
			<td><input type="text" name="title" /></td>
		</tr>
		<tr>
			<td>Description</td>
			<td><input type="text" name="description" /></td>
		</tr>
	</table>
	<button type="submit">Create Thread</button>
</form>
`

func (h *Handler) ThreadsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(threadsCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) ThreadsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.store.CreateThread(&miji.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *Handler) ThreadsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.store.DeleteThread(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

const threadsUpdateHTML = `
<h1>Update Thread</h1>
<form action="/threads/update" method="POST">
	<table>
		<input type="hidden" name="id" value="{{.ID}}" />
		<tr>
			<td>Title</td>
			<td><input type="text" name="title" value="{{.Title}}" /></td>
		</tr>
		<tr>
			<td>Description</td>
			<td><input type="text" name="description" value="{{.Description}}" /></td>
		</tr>
	</table>
	<button type="submit">Update Thread</button>
</form>
`

func (h *Handler) ThreadsUpdatePage() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(threadsUpdateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpl.Execute(w, map[string]string{
			"ID":          idStr, // this instead of id.String(); change my mind
			"Title":       p.Title,
			"Description": p.Description,
		})
	}
}

func (h *Handler) ThreadsUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.FormValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.store.UpdateThread(&miji.Thread{
			ID:          id,
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
