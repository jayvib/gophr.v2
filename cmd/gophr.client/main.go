package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jayvib/golog"
)

var (
	port  string
	host  string
	debug bool
)

var (
	templates      *template.Template
	layoutTemplate *template.Template
)

var funcs = template.FuncMap{
	"yield": func() (string, error) {
		return "", nil
	},
	"navbar": func() (string, error) {
		return "", nil
	},
}

func init() {
	flag.StringVar(&port, "port", "8080", "Port")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&host, "host", "localhost", "Host")
	flag.Parse()

	if debug {
		golog.Warning("DEBUGGING MODE!")
		golog.SetLevel(golog.DebugLevel)
	}

	templates = template.Must(
		template.ParseGlob("templates/**/*.html"))

	layoutTemplate = template.Must(
		template.New("layout.html").
			Funcs(funcs).
			ParseFiles("templates/layout.html"))
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	router.HandleFunc("/", HomePage).Methods(http.MethodGet)
	router.HandleFunc("/signup", Signup).Methods(http.MethodGet)
	router.HandleFunc("/login", Login).Methods(http.MethodGet)
	addr := fmt.Sprintf("%s:%s", host, port)
	golog.Info("Listening to", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
