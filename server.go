package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var userID = 0

func StartWebServer() {
	var lobby = NewLobby()
	var comm = NewCommunicator(lobby)

	http.HandleFunc("/app.html", func(w http.ResponseWriter, r *http.Request) {
		user := ""
		for _, c := range r.Cookies() {
			if c.Name == "user_id" {
				user = c.Value
			}
		}
		if user == "" {
			user = fmt.Sprintf("%d", userID)
			userID++
		}
		c := http.Cookie{
			Name: "user_id",
			Value: user,
		}
		http.SetCookie(w, &c)
		http.ServeFile(w, r, filepath.Join("./www/app.html"))
	})

	http.HandleFunc("/ws", comm.serveWs)

	http.HandleFunc("/create_game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			gs := GameSetting{NumberOfPlayers: 4, Map: 1, TurnTimeOut: false}
			gameID, _ := lobby.CreateGame(gs)
			w.Write([]byte(fmt.Sprintf(`{"g":"%d"}`, gameID)))
		}
	})

	http.HandleFunc("/update_settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			gameID := r.FormValue("g")
			g, err := strconv.Atoi(gameID)
			if err != nil {
				return
			}
			game := lobby.GetGame(g)
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
			g, err := strconv.Atoi(gameID)
			if err != nil {
				return
			}
			game := lobby.GetGame(g)
			err = game.Start()
			if err != nil {
				w.Write([]byte(fmt.Sprint(`{"ok":true}`)))
			} else {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}
		}
	})

	http.HandleFunc("/join_game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			g := r.FormValue("g")
			gameID, err := strconv.Atoi(g)
			if err != nil {
				return
			}
			userID := r.FormValue("user_id")
			lobby.AddUserToGame(gameID, userID)
		}
	})

	http.HandleFunc("/roll_dice", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			gameID := r.FormValue("g")
			g, err := strconv.Atoi(gameID)
			if err != nil {
				return
			}
			game := lobby.GetGame(g)
			err = game.RollDice()
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
			g, err := strconv.Atoi(gameID)
			if err != nil {
				return
			}
			game := lobby.GetGame(g)
			err = game.BankTrade([2]int{gCardType, gCardCount}, wantCardType)
			if err != nil {
				w.Write([]byte(fmt.Sprint(`{"ok":true}`)))
			} else {
				w.Write([]byte(fmt.Sprint(`{"ok":false}`)))
			}
		}
	})

	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			gameID := r.FormValue("g")
			g, err := strconv.Atoi(gameID)
			if err != nil {
				return
			}
			game := lobby.GetGame(g)
			w.Write([]byte(game.Board()))
		}
	})

	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", fs)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
