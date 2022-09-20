package app

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
