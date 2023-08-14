package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorlog *log.Logger
	infolog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "Http Network Address")
	flag.Parse()
	loginfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	app := application{
		errorlog: errlog,
		infolog:  loginfo,
	}

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errlog,
		Handler:  mux,
	}

	loginfo.Printf("Starting Server On port : %s \n", *addr)
	err := srv.ListenAndServe()
	errlog.Fatal(err)

}
