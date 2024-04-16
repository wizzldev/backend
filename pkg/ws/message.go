package ws

type Map map[string]interface{}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
