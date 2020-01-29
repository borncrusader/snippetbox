package main

import (
	"net/http"
	"strconv"

	"srinathkrishna.in/snippetbox/pkg/models"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippets: s,
	}

	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showSnippetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	data := &templateData{
		Snippet: s,
	}

	app.render(w, r, "show.page.tmpl", data)
}

func (app *application) createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "foo bar"
	expires := "7"
	_, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/snippet", http.StatusSeeOther)
}
