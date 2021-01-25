package webserver

import (
	"catans/board"
	"fmt"
	"log"
	"net/http"
)

func StartWebServer() {
	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, board.CatNodes(0))
	})

	http.HandleFunc("/intersections", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, board.CatIntersections(0))
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
