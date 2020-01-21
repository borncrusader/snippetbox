package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type application struct {
	addr     *string
	dsn      *string
	mux      *http.ServeMux
	server   *http.Server
	errorLog *log.Logger
	infoLog  *log.Logger
}

func (app *application) parseArgs() {
	app.addr = flag.String("addr", ":4000", "HTTP Network Address")
	app.dsn = flag.String("dsn", "user=web password=password host=localhost port=5432 database=snippetbox sslmode=disable", "PGX DSN")
	flag.Parse()
}

func (app *application) setupConfig() {
	app.errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	app.infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.LUTC)
}

func (app *application) openDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", *app.dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) createServer() {
	app.mux = http.NewServeMux()
	app.server = &http.Server{
		Addr:     *app.addr,
		ErrorLog: app.errorLog,
		Handler:  app.mux,
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

	db, err := app.openDB()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	app.createServer()
	app.registerRoutes()

	err = app.startServer()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
