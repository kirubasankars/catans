package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func StartWebServer() {
	lobby := NewLobby()

	http.HandleFunc("/create_game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			gs := GameSetting{NumberOfPlayers: 4, Map: 0}
			gameID, _ := lobby.CreateGame(gs)
			w.Write([]byte(fmt.Sprintf(`{"g":"%s"}`, gameID)))
		}
	})

	http.HandleFunc("/update_settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			gameID := r.FormValue("g")
			game := lobby.GetGame(gameID)
			gs := game.context.GameSetting

			m := r.FormValue("map")
			mi, _ := strconv.Atoi(m)
			gs.Map = mi

			game.UpdateGameSetting(gs)
			w.Write([]byte(fmt.Sprint(`{"ok":true}`)))
		}
	})

	http.HandleFunc("/start_game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			gameID := r.FormValue("g")
			game := lobby.GetGame(gameID)
			err := game.Start()
			if err != nil {
				w.Write([]byte(fmt.Sprint(`{"ok":true}`)))
			} else {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}
		}
	})

	http.HandleFunc("/roll_dice", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			gameID := r.FormValue("g")
			game := lobby.GetGame(gameID)
			err := game.RollDice()
			if err != nil {
				w.Write([]byte(fmt.Sprint(`{"ok":true}`)))
			} else {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}
		}
	})

	http.HandleFunc("/bank_trade", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()

			gives := r.FormValue("gives")
			s := strings.Split(gives, ",")
			if len(s) != 2 {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}
			gCardType, err := strconv.Atoi(s[0])
			if err != nil || gCardType < 0 || gCardType > 4 {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}
			gCardCount, err := strconv.Atoi(s[1])
			if err != nil || gCardCount > 4 || gCardCount < 2 {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}

			want := r.FormValue("want")
			wantCardType, err := strconv.Atoi(want)
			if err != nil || gCardType < 0 || gCardType > 4 {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}

			gameID := r.FormValue("g")
			game := lobby.GetGame(gameID)
			err = game.BankTrade([2]int{gCardType, gCardCount}, wantCardType)
			if err != nil {
				w.Write([]byte(fmt.Sprint(`{"ok":true}`)))
			} else {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}
		}
	})

	http.HandleFunc("/board_ui", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			gameID := r.FormValue("g")
			game := lobby.GetGame(gameID)
			w.Write([]byte(game.UI()))
		}
	})

	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", fs)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
