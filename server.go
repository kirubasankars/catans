package main

import (
	"fmt"
	"log"
	"net/http"
)

func StartWebServer() {
	lobby := NewLobby()

	http.HandleFunc("/create_game", func(w http.ResponseWriter, r *http.Request) {
		gs := GameSetting{NumberOfPlayers: 2}
		gameID, _ := lobby.CreateGame(gs)
		w.Write([]byte(gameID))
	})

	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		game := lobby.GetGame("GAME#1")
		w.Write([]byte(game.UI()))
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
