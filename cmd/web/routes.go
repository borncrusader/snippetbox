package main

import "net/http"

func (app *application) registerRoutes() http.Handler {
	app.mux.Get("/", http.HandlerFunc(app.homeHandler))
	app.mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetFormHandler))
	app.mux.Post("/snippet/create", http.HandlerFunc(app.createSnippetHandler))
	// this should be lower since it's less specific than
	app.mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippetHandler))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	app.mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	return app.recoverPanic(
		app.recoverTimeouts(
			app.logRequest(
				app.secureHeaders(
					app.mux,
				),
			),
		),
	)
}
