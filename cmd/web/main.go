package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/devruhulamin/go-snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorlog      *log.Logger
	infolog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "Http Network Address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()
	loginfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(*dsn)
	if err != nil {
		errlog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		errlog.Fatal(err)
	}

	app := application{
		errorlog: errlog,
		infolog:  loginfo,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errlog,
		Handler:  app.routes(),
	}

	loginfo.Printf("Starting Server On port : %s \n", *addr)
	err = srv.ListenAndServe()
	errlog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
