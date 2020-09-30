package web

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gocs/miji"
	"github.com/google/uuid"
)

const threadsListHTML = `
<h1>Threads</h1>
<form action="/threads/new" method="GET">
	<button type="submit">Create New</button>
</form>
<dl>{{range .Threads}}
	<dt><strong>{{.Title}}</strong></dt>
	<dd>{{.Description}}</dd>
	<dd><form action="/threads/{{.ID}}" method="GET">
			<button type="submit">Show</button>
		</form><form action="/threads/{{.ID}}/delete" method="POST">
			<button type="submit">Delete</button>
		</form>
	</dd>
{{end}}</dl>
`
 
type ThreadHandler struct {
	store miji.Store
}

func (h *ThreadHandler) List() http.HandlerFunc {
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

func (h *ThreadHandler) Create() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(threadsCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *ThreadHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		t := &miji.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}

		if err := h.store.CreateThread(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads/"+t.ID.String(), http.StatusFound)
	}
}

func (h *ThreadHandler) Delete() http.HandlerFunc {
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

const threadsShowHTML = `
<h1>Show Thread</h1>
<form action="/threads/update" method="POST">
	<table>
		<input type="hidden" name="id" value="{{.Thread.ID}}" />
		<tr>
			<td>Title</td>
			<td><input type="text" name="title" value="{{.Thread.Title}}" /></td>
		</tr>
		<tr>
			<td>Description</td>
			<td><input type="text" name="description" value="{{.Thread.Description}}" /></td>
		</tr>
	</table>
	<button type="submit">Update Thread</button>
</form>
<form action="/threads" method="GET">
	<button type="submit">Cancel</button>
</form>
<form action="/threads/{{.Thread.ID}}/new" method="GET">
	<button type="submit">New Post</button>
</form>
<h1>Show Posts</h1>
<dl>{{range .Posts}}
	<dt><strong>{{.Title}}</strong> - {{.Votes}}</dt>
	<dd>{{.Content}}</dd>
	<dd><form action="/threads/{{.ThreadID}}/{{.ID}}" method="GET">
			<button type="submit">Show</button>
		</form><form action="/threads/{{.ThreadID}}/{{.ID}}/delete" method="POST">
			<button type="submit">Delete</button>
		</form>
	</dd>
{{end}}</dl>
`

func (h *ThreadHandler) Show() http.HandlerFunc {
	type data struct {
		Thread miji.Thread
		Posts   []miji.Post
	}

	tmpl := template.Must(template.New("").Parse(threadsShowHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ps, err := h.store.PostsByThread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpl.Execute(w, data{
			Thread: t,
			Posts:   ps,
		})
	}
}

func (h *ThreadHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.FormValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		title := r.FormValue("title")
		description := r.FormValue("description")

		t := &miji.Thread{
			ID:          id,
			Title:       title,
			Description: description,
		}

		if err := h.store.UpdateThread(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads/"+t.ID.String(), http.StatusFound)
	}
}
