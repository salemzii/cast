package broadcast

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/salemzii/cast.git/db"
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
}

func (c *Client) ReadDb() {
	ticker := time.NewTicker(1 * time.Second)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	var rides []db.Ride
	var err error
	for {
		select {
		case <-ticker.C:
			rides, err = db.RideRespository.All()
			if err != nil {
				log.Println(err)
			}
			fmt.Println("here", len(rides))
			for _, ride := range rides {
				fmt.Println(ride)
				rideMessage := []byte(fmt.Sprintf("%s:%s", ride.Clientid, ride.RideId))
				c.hub.ride <- rideMessage

			}
		}
	}
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

	client := &Client{hub: hub, clientId: r.URL.RawQuery, driverId: r.URL.RawQuery, conn: conn, send: make(chan []byte, 256)}
	if strings.HasPrefix(client.clientId, "driverId") {
		client.clientId = ""
		client.Type = "driver"
	} else if strings.HasPrefix(client.driverId, "clientId") {
		client.driverId = ""
		client.Type = "client"
	}

	client.hub.register <- client

	fmt.Println(client.Type)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	//go client.ReadDb()
	go client.writePump()
	go client.readPump()
}
