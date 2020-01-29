package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"srinathkrishna.in/snippetbox/pkg/models/pgsql"
)

type application struct {
	addr              *string
	dsn               *string
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	defaultTimeout    time.Duration
	mux               *http.ServeMux
	server            *http.Server
	errorLog          *log.Logger
	infoLog           *log.Logger
	db                *sql.DB
	snippets          *pgsql.SnippetModel
	templateCache     map[string]*template.Template
}

func (app *application) parseArgs() {
	app.addr = flag.String("addr", ":4000", "HTTP Network Address")
	app.dsn = flag.String("dsn", "user=web password=password host=localhost port=5432 database=snippetbox sslmode=disable", "PGX DSN")

	// TODO: these need to be parsed from args
	app.readTimeout = 5 * time.Second
	app.readHeaderTimeout = 5 * time.Second
	app.writeTimeout = 5 * time.Second
	app.idleTimeout = 10 * time.Second
	app.defaultTimeout = 2 * time.Second

	flag.Parse()
}

func (app *application) setupConfig() {
	app.errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	app.infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.LUTC)
}

func (app *application) setupDB() error {
	var err error

	app.db, err = sql.Open("pgx", *app.dsn)
	if err != nil {
		return err
	}

	if err = app.db.Ping(); err != nil {
		return err
	}

	app.snippets = &pgsql.SnippetModel{DB: app.db}

	return nil
}

func (app *application) primeCaches() {
	cache, err := newTemplateCache("./ui/html")
	if err != nil {
		app.errorLog.Fatal(err)
	}

	app.templateCache = cache
}

func (app *application) createServer() {
	app.mux = http.NewServeMux()
	app.server = &http.Server{
		Addr:              *app.addr,
		ErrorLog:          app.errorLog,
		ReadTimeout:       app.readTimeout,
		ReadHeaderTimeout: app.readHeaderTimeout,
		WriteTimeout:      app.writeTimeout,
		IdleTimeout:       app.idleTimeout,
		// TODO: have a better timeout page
		Handler: http.TimeoutHandler(app.mux, app.defaultTimeout, "Timeout!"),
	}
}

func (app *application) startServer() error {
	app.infoLog.Printf("Starting server on %v", *app.addr)

	return http.ListenAndServe(*app.addr, app.mux)
}

func main() {
	app := &application{}

	app.parseArgs()
	app.setupConfig()

	err := app.setupDB()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer app.db.Close()

	app.primeCaches()

	app.createServer()
	app.registerRoutes()

	err = app.startServer()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
