package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "server/home.html")
}

func Serve() {
	hub := newHub()
	go hub.run()
	ro := mux.NewRouter() // Create a mux instance
	ro.HandleFunc("/", serveHome)
	ro.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(":8088", ro)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
