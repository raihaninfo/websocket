package main

import (
	"log"
	"net/http"
	"text/template"
)

func (a *application) welcome(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		log.Println(err)
	}
	tpl.Execute(w, nil)
}

//static file server
