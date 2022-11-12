package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

func (a *application) welcome(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		log.Println(err)
	}
	tpl.Execute(w, nil)
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

func (a *application) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected to Endpoint")
	var response WsJsonResponse
	response.Message = `<em>Client Connected</em>`
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
}
