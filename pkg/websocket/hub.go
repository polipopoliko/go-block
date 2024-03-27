package ws

import "log"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Register(cl *Client) {
	h.register <- cl
}

func (h *Hub) Unregister(cl *Client) {
	h.unregister <- cl
}

func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

func (h *Hub) PrintClients() {
	log.Println(h.clients)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.buffer)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.buffer <- message:
				default:
					close(client.buffer)
					delete(h.clients, client)
				}
			}
		}
	}
}
