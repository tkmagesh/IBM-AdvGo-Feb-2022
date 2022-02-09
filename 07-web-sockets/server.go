package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsConnections []*websocket.Conn

func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	wsConnections = append(wsConnections, ws)
	go reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msg := string(p)
		log.Println(msg)
		broadcast(msg)
	}
}

func broadcast(msg string) {
	for _, conn := range wsConnections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func cleanup() {
	ticker := time.Tick(5 * time.Second)
	for range ticker {
		fmt.Println("connects # : ", len(wsConnections))
	}

}

func main() {
	go cleanup()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/chat", wsEndPoint)
	http.ListenAndServe(":8080", nil)
}
