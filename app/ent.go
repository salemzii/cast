package app

import (
	"github.com/gorilla/websocket"
)

type App struct {
	AppId   string `json:"appid"`
	Devices []Device
}

type Device struct {
	DeviceId string `json:"deviceid"`
}

type Message struct {
	AppId    string `json:"appid"`
	ClientId string `json:"clientid"`
	Data     string `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
