package main

import (
	"flag"
	"log"
	"net/http"

	broadcast "github.com/salemzii/cast.git/broadcast"
	"github.com/salemzii/cast.git/db"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	flag.Parse()

	hub := broadcast.NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/home", serveHome2)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		broadcast.ServeWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func DataEntry() {
	df := []db.Ride{
		{RideId: "01234", Clientid: "1002", Status: "initiated"},
		{RideId: "8393", Clientid: "1003", Status: "initiated"},
		{RideId: "3421", Clientid: "1004", Status: "processed"},
		{RideId: "2540", Clientid: "1005", Status: "initiated"},
		{RideId: "0004", Clientid: "1006", Status: "initiated"},
	}
	for _, v := range df {
		_, err := db.RideRespository.Create(v)
		if err != nil {
			log.Println(err)
		}
	}
}

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
	http.ServeFile(w, r, "home.html")
}
func serveHome2(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/home" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home2.html")
}
