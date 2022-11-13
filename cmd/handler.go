package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"text/template"

	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[websocketConnection]string)

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

type websocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action        string   `json:"action"`
	Message       string   `json:"message"`
	MessageType   string   `json:"message_type"`
	ConnectedUser []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	UserName string              `json:"username"`
	Message  string              `json:"message"`
	Conn     websocketConnection `json:"_"`
}

func (a *application) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected to Endpoint")
	var response WsJsonResponse
	response.Message = `<em>Client Connected</em>`

	conn := websocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *websocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			// get a list of all user
			clients[e.Conn] = e.UserName
			users := getUserList()
			// fmt.Println(users)
			response.Action = "list_users"
			response.ConnectedUser = users
			broadcastToAll(response)

		}

		// response.Action = "Got hear"
		// response.Message = fmt.Sprintf("Some message, and action was %s", e.Action)
		// broadcastToAll(response)

	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList
}

func broadcastToAll(res WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(res)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}
