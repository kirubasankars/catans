package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

type Communicator struct {
	lobby *GameLobby
	upgrader websocket.Upgrader
}

func (communicator *Communicator) addUser(ID string, con *websocket.Conn) {
	user := NewUser(ID, "Dev`")
	user.Status = 1
	user.con = con
	communicator.lobby.users[ID] = *user
}

func (communicator *Communicator) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := communicator.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	//read id from cookie
	user := ""
	for _, c := range r.Cookies() {
		if c.Name == "user_id" {
			user = c.Value
		}
	}
	communicator.addUser(user, ws)
}

func (communicator *Communicator) ping() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, user := range communicator.lobby.users {
				if err := user.con.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
					log.Println("ping:", err)
					user.Status = 0
					user.con.Close()
					user.con = nil
				}
			}
		}
	}
}

func NewCommunicator(lobby *GameLobby) *Communicator {
	comm := new(Communicator)
	comm.lobby = lobby
	go comm.ping()
	return comm
}