package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/salemzii/cast.git/app"
	"github.com/salemzii/cast.git/broadcast"
)

func main() {

	flag.Parse()
	port := os.Getenv("PORT")
	hub := broadcast.NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/publish", Publish)
	http.HandleFunc("/taskify/notify", PublishTaskifyNotification)
	http.HandleFunc("/taskify/chat", PublishTaskifyChat)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		broadcast.ServeWs(hub, w, r)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
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

func Publish(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/publish" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg app.Message
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(data, &msg)

	if err != nil {
		log.Println(err)
	}
	broadcast.MessageQueue <- msg
}

func PublishTaskifyNotification(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/taskify/notify" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg app.Message
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(data, &msg)

	if err != nil {
		log.Println(err)
	}
	broadcast.MessageQueue <- msg
}

func PublishTaskifyChat(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/taskify/chat" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg app.Message
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(data, &msg)

	if err != nil {
		log.Println(err)
	}
	broadcast.MessageQueue <- msg
}
