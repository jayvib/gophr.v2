package main

import (
	"bytes"
	"github.com/jayvib/golog"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, "index/navbar", nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, "users/signup", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, "sessions/login", nil)
}

func renderPage(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {

	// create a custom funcs that has a yield function.
	// The function will return an template.HTML type and an error
	f := template.FuncMap{
		"navbar": func() (template.HTML, error) {
			var buff bytes.Buffer
			err := templates.ExecuteTemplate(&buff, "index/navbar", data)
			return template.HTML(buff.String()), err
		},
	}

	if name != "" {
		f["yield"] = func() (template.HTML, error) {
			var buff bytes.Buffer
			err := templates.ExecuteTemplate(&buff, name, data)
			return template.HTML(buff.String()), err
		}
	}

	// clone a layout
	clonedLayout, _ := layoutTemplate.Clone()
	clonedLayout.Funcs(f)

	// execute the layout
	err := clonedLayout.Execute(w, nil)
	if err != nil {
		err = templates.ExecuteTemplate(w, "other/error", map[string]interface{}{
			"Error": err.Error(),
		})
		if err != nil {
			golog.Error(err)
		}
	}
}
