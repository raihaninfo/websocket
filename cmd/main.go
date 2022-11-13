package main

import "log"

type application struct {
	AppName string
	Server  Server
	Debug   bool
}

type Server struct {
	Host string
	Port string
	Url  string
}

func main() {
	app := application{
		AppName: "My App",
		Server: Server{
			Host: "localhost",
			Port: "8080",
			Url:  "localhost:8080",
		},
		Debug: true,
	}
	log.Printf("Listening server on http://localhost:%v", app.Server.Port)
	go ListenToWsChannel()
	app.ListenAndServe()
}
