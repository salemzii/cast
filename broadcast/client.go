package broadcast

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/salemzii/cast.git/app"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	MessageQueue = make(chan app.Message, 300)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	clientId string
	driverId string
	Type     string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	sendMsg chan app.Message
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("NEW CLIENT CONNECTION: ")
	fmt.Println(r.URL.RawQuery)
	fmt.Println("--------------------------")

	client := &Client{hub: hub, clientId: r.URL.RawQuery, driverId: r.URL.RawQuery,
		conn: conn, send: make(chan []byte, 256)}

	if strings.HasPrefix(client.clientId, "driverId") {
		client.clientId = ""
		client.Type = "driver"
	} else if strings.HasPrefix(client.driverId, "clientId") {
		client.driverId = ""
		client.Type = "client"
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	//go client.ReadDb()
	go client.writePump()
	go client.readPump()
}
