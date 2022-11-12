package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *application) router() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", a.welcome).Methods("GET")
	mux.HandleFunc("/ws", a.WsEndpoint).Methods("GET")
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", mux)

	return mux
}
