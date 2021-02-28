package main

import (
	"github.com/gorilla/websocket"
)

type User struct {
	ID                   string
	Name                 string
	Status               int
	LastPublishedEventID int

	con *websocket.Conn
}

func NewUser(ID string, Name string) *User {
	user := new(User)
	user.ID = ID
	user.Name = Name
	return user
}
