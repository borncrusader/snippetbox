package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	addr     *string
	mux      *http.ServeMux
	server   *http.Server
	errorLog *log.Logger
	infoLog  *log.Logger
}

func (app *application) parseArgs() {
	app.addr = flag.String("addr", ":4000", "HTTP Network Address")

	flag.Parse()
}

func (app *application) setupConfig() {
	app.errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	app.infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.LUTC)
}

func (app *application) createServer() {
	app.mux = http.NewServeMux()
	app.server = &http.Server{
		Addr:     *app.addr,
		ErrorLog: app.errorLog,
		Handler:  app.mux,
	}
}

func (app *application) startServer() {
	app.infoLog.Printf("Starting server on %v", *app.addr)

	err := http.ListenAndServe(*app.addr, app.mux)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func main() {
	app := &application{}

	app.parseArgs()
	app.setupConfig()

	app.createServer()
	app.registerRoutes()

	app.startServer()
}
