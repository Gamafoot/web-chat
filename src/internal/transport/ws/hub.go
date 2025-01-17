package ws

import (
	"fmt"
)

type Clients map[int]*Client

type Rooms map[string]*Room

type Hub struct {
	Rooms      Rooms
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(Rooms),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomName]; ok {
				r := h.Rooms[client.RoomName]

				if _, ok := r.Clients[client.Id]; !ok {
					r.Clients[client.Id] = client
				}
			}

		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomName]; ok {
				if _, ok := h.Rooms[client.RoomName].Clients[client.Id]; ok {
					h.Broadcast <- &Message{
						Type:     "event",
						RoomName: client.RoomName,
						Color:    client.Color,
						Content:  "user left the chat",
					}

					h.Rooms[client.RoomName].ColorDispanser.Reset(client.Color)

					delete(h.Rooms[client.RoomName].Clients, client.Id)

					if len(h.Rooms[client.RoomName].Clients) == 0 {
						delete(h.Rooms, client.RoomName)
					}

					close(client.Message)
				}
			}

		case msg := <-h.Broadcast:
			if _, ok := h.Rooms[msg.RoomName]; ok {
				for _, cl := range h.Rooms[msg.RoomName].Clients {
					cl.Message <- msg
				}
			}
		}
	}
}

func (h *Hub) GetRoomByName(name string) (*Room, error) {
	for _, r := range h.Rooms {
		if r.Name == name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("room (%s) not exists", name)
}
