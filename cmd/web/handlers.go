package main

import (
	"fmt"
	"net/http"
	"strconv"

	"srinathkrishna.in/snippetbox/pkg/models"
)

func (app *application) handleHomeGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := app.snippets.Latest(10)
		if err != nil {
			app.serverError(w, err)
			return
		}

		app.render(w, r, "home.page.tmpl", &templateData{
			Snippets: s,
		})
	}
}

func (app *application) handleSnippetGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		s, err := app.snippets.Get(id)
		if err == models.ErrNoRecord {
			app.notFound(w)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		app.render(w, r, "show.page.tmpl", &templateData{
			Snippet: s,
		})
	}
}

func (app *application) handleSnippetCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// limit size of the body, ParseForm() would fail if request body was more than 10M
		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")
		expires := r.PostForm.Get("expires")

		id, err := app.snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	}
}

func (app *application) handleSnippetCreateForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "create.page.tmpl", nil)
	}
}
