package webserver

import (
	"catans/board"
	"fmt"
	"log"
	"net/http"
)

func StartWebServer() {
	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.FormValue("name")
		fmt.Fprintf(w, board.CatNodes(name))
	})

	http.HandleFunc("/intersections", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.FormValue("board")
		fmt.Fprintf(w, board.CatIntersections(name))
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
