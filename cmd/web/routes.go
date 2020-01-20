package main

import "net/http"

func (app *application) registerRoutes() {
	app.mux.HandleFunc("/", app.homeHandler)
	app.mux.HandleFunc("/snippet", app.showSnippetHandler)
	app.mux.HandleFunc("/snippet/create", app.createSnippetHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	app.mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}
