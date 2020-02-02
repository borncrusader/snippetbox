package main

import (
	"fmt"
	"net/http"
	"strconv"

	"srinathkrishna.in/snippetbox/pkg/forms"
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

// handleSnippetCreate handles the POST
func (app *application) handleSnippetCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// limit size of the body, ParseForm() would fail if request body was more than 10M
		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		form := forms.New(r.Form)
		form.Required("title", "content", "expires")
		form.MaxLength("title", 100)
		form.PermittedValues("expires", "365", "7", "1")

		if !form.Valid() {
			app.render(w, r, "create.page.tmpl", &templateData{
				Form: form,
			})
			return
		}

		id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	}
}

func (app *application) handleSnippetCreateForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}
}

func (app *application) handleUserSignupForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}
}

// handleUserSignup handles the POST
func (app *application) handleUserSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// limit size of the body, ParseForm() would fail if request body was more than 10M
		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		form := forms.New(r.Form)
		form.Required("name", "email", "password")
		form.MaxLength("name", 255)
		form.MaxLength("email", 255)
		form.MinMaxLength("password", 8, 32)

		if !form.Valid() {
			app.render(w, r, "signup.page.tmpl", &templateData{
				Form: form,
			})
			return
		}

		//id, err := app.users.Insert(form.Get("name"), form.Get("email"),
		//							form.Get(""))

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *application) handleUserLoginForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (app *application) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (app *application) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
