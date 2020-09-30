package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gocs/miji"
	"github.com/google/uuid"
)

type PostHandler struct {
	store miji.Store
}

const postsCreateHTML = `
<h1>New Post</h1>
<form action="/threads/{{.Thread.ID}}" method="POST">
	<table>
		<tr>
			<td>Title</td>
			<td><input type="text" name="title" /></td>
		</tr>
		<tr>
			<td>Content</td>
			<td><input type="text" name="content" /></td>
		</tr>
	</table>
	<button type="submit">Create Post</button>
</form>
<form action="/threads/{{.Thread.ID}}" method="GET">
	<button type="submit">Cancel</button>
</form>
`

func (h *PostHandler) Create() http.HandlerFunc {
	type data struct {
		Thread miji.Thread
	}

	tmpl := template.Must(template.New("").Parse(postsCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			Thread: t,
		})
	}
}

func (h *PostHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		content := r.FormValue("content")

		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := &miji.Post{
			ID:       uuid.New(),
			ThreadID: t.ID,
			Title:    title,
			Content:  content,
		}
		if err := h.store.CreatePost(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/threads/"+t.ID.String()+"/"+p.ID.String(), http.StatusFound)
	}
}

func (h *PostHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "postID")
		threadIDStr := chi.URLParam(r, "threadID")

		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.store.DeletePost(postID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads/"+threadID.String(), http.StatusFound)
	}
}

const postsShowHTML = `
<h1>Show Post</h1>
<form action="/threads/{{.Thread.ID}}/update" method="POST">
	<table>
		<input type="hidden" name="id" value="{{.Post.ID}}" />
		<tr>
			<td>Title</td>
			<td><input type="text" name="title" value="{{.Post.Title}}" /></td>
		</tr>
		<tr>
			<td>Content</td>
			<td><input type="text" name="content" value="{{.Post.Content}}" /></td>
		</tr>
		<tr>
			<td>Votes</td>
			<td>
				<input type="hidden" name="votes" value="{{.Post.Votes}}" />
				<span>{{.Post.Votes}}</span>
			</td>
		</tr>
	</table>
	<button type="submit">Update Post</button>
</form>
<dl>{{range .Comments}}
	<dd>{{.Content}} - {{.Votes}}</dd>
	<dd><form action="/threads/{{.ID}}" method="GET">
			<button type="submit">Show</button>
		</form><form action="/threads/{{.ID}}/delete" method="POST">
			<button type="submit">Delete</button>
		</form>
	</dd>
{{end}}</dl>
<form action="/threads/{{.Thread.ID}}/{{.Post.ID}}/vote?dir=up" method="POST">
	<button type="submit">Upvote</button>
</form>
<form action="/threads/{{.Thread.ID}}/{{.Post.ID}}/vote?dir=down" method="POST">
	<button type="submit">Downvote</button>
</form>
<form action="/threads/{{.Thread.ID}}" method="GET">
	<button type="submit">Cancel</button>
</form>
`

func (h *PostHandler) Show() http.HandlerFunc {
	type data struct {
		Thread   miji.Thread
		Post     miji.Post
		Comments []miji.Comment
	}

	tmpl := template.Must(template.New("").Parse(postsShowHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "postID")
		threadIDStr := chi.URLParam(r, "threadID")

		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.Post(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cc, err := h.store.CommentsByPost(p.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := h.store.Thread(threadID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			Thread:   t,
			Post:     p,
			Comments: cc,
		})
	}
}

func (h *PostHandler) Vote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "postID")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.Post(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			p.Votes++
		} else if dir == "down" {
			p.Votes--
		}

		if err := h.store.UpdatePost(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (h *PostHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")

		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		postIDStr := r.FormValue("id")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		votesStr := r.FormValue("votes")
		votes, err := strconv.Atoi(votesStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := &miji.Post{
			ID:       postID,
			ThreadID: threadID,
			Title:    title,
			Content:  content,
			Votes:    votes,
		}
		

		if err := h.store.UpdatePost(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads/"+threadIDStr+"/"+postIDStr, http.StatusFound)
	}
}
