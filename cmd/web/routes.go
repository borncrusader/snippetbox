package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) registerRoutes() http.Handler {
	defaultChain := alice.New(app.recoverPanic, app.recoverTimeouts,
		app.logRequest, app.secureHeaders)

	dynamicChain := alice.New(app.session.Enable)

	app.router = pat.New()

	app.router.Get("/", dynamicChain.ThenFunc(app.handleHomeGet()))
	app.router.Get("/snippet/create", dynamicChain.ThenFunc(app.handleSnippetCreateForm()))
	app.router.Post("/snippet/create", dynamicChain.ThenFunc(app.handleSnippetCreate()))
	// this should be lower since it's less specific
	app.router.Get("/snippet/:id", dynamicChain.ThenFunc(app.handleSnippetGet()))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	app.router.Get("/static/", http.StripPrefix("/static/", fileServer))

	return defaultChain.Then(app.router)
}
