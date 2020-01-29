package main

import "net/http"

func (app *application) registerRoutes() http.Handler {
	app.mux.HandleFunc("/", app.homeHandler)
	app.mux.HandleFunc("/snippet", app.showSnippetHandler)
	app.mux.HandleFunc("/snippet/create", app.createSnippetHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	app.mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// TODO: have a better timeout page
	return app.recoverPanic(
		http.TimeoutHandler(
			app.logRequest(
				app.secureHeaders(
					app.mux,
				),
			), app.defaultTimeout, "Timeout!",
		),
	)
}
