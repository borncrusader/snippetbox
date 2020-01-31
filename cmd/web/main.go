package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bmizerany/pat"
	"github.com/golangcollege/sessions"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"srinathkrishna.in/snippetbox/pkg/models/pgsql"
)

type application struct {
	addr              *string
	db                *sql.DB
	defaultTimeout    time.Duration
	dsn               *string
	errorLog          *log.Logger
	idleTimeout       time.Duration
	infoLog           *log.Logger
	readHeaderTimeout time.Duration
	readTimeout       time.Duration
	router            *pat.PatternServeMux
	secret            *string
	server            *http.Server
	session           *sessions.Session
	snippets          *pgsql.SnippetModel
	templateCache     map[string]*template.Template
	writeTimeout      time.Duration
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

	app.secret = flag.String("secret", "872ADm1srQkgqwy2E39h33OqTKLwnYUf", "Secret Key")

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
	app.session = sessions.New([]byte(*app.secret))
	app.session.Lifetime = 12 * time.Hour
	app.server = &http.Server{
		Addr:              *app.addr,
		ErrorLog:          app.errorLog,
		ReadTimeout:       app.readTimeout,
		ReadHeaderTimeout: app.readHeaderTimeout,
		WriteTimeout:      app.writeTimeout,
		IdleTimeout:       app.idleTimeout,
		Handler:           app.registerRoutes(),
	}
}

func (app *application) startServer() error {
	app.infoLog.Printf("Starting server on %v", *app.addr)

	return app.server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
}

func run(app *application) error {
	app.parseArgs()

	app.setupConfig()

	err := app.setupDB()
	if err != nil {
		return errors.Wrap(err, "setup database")
	}
	defer app.db.Close()

	app.primeCaches()

	app.createServer()

	err = app.startServer()
	if err != nil {
		return errors.Wrap(err, "server start")
	}

	return nil
}

func main() {
	app := &application{}

	if err := run(app); err != nil {
		app.errorLog.Fatal(err)
		os.Exit(1)
	}
}
