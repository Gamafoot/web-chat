package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	Id       int    `json:"id"`
	Username string `json:"username"`
	RoomName string `json:"roomName"`
	Color    string `json:"color"`
}

type Message struct {
	Type     string `json:"type"`
	UserId   int    `json:"userId"`
	Color    string `json:"color"`
	RoomName string `json:"roomName"`
	Content  string `json:"content"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		msg, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(msg)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %+v\n", err)
			}
			break
		}

		msg := &Message{
			UserId:   c.Id,
			RoomName: c.RoomName,
			Content:  string(m),
			Color:    c.Color,
		}

		hub.Broadcast <- msg
	}
}
