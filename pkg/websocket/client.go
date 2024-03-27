package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub    *Hub
	buffer chan []byte
	Conn   *websocket.Conn
}

func NewWebsocketClient(readBufferSize, writeBufferSize int, hub *Hub, w http.ResponseWriter, r *http.Request) *Client {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader err", err)
		return nil
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	return &Client{
		hub:    hub,
		buffer: make(chan []byte, 256),
		Conn:   conn,
	}
}

func (cl *Client) Close() {
	if cl == nil {
		return
	}

	if cl.Conn == nil {
		return
	}

	cl.Conn.Close()
}
