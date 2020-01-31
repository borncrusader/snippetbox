package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) registerRoutes() http.Handler {
	chain := alice.New(app.recoverPanic, app.recoverTimeouts,
		app.logRequest, app.secureHeaders)

	app.router = pat.New()

	app.router.Get("/", app.handleHomeGet())
	app.router.Get("/snippet/create", app.handleSnippetGet())
	app.router.Post("/snippet/create", app.handleSnippetCreateForm())
	// this should be lower since it's less specific
	app.router.Get("/snippet/:id", app.handleSnippetGet())

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	app.router.Get("/static/", http.StripPrefix("/static/", fileServer))

	return chain.Then(app.router)
}
