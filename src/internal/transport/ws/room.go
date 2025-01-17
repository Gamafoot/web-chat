package ws

import (
	"net/http"
	colordispanser "root/pkg/color_dispanser"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

type Room struct {
	Name           string
	Clients        Clients
	ColorDispanser *colordispanser.Dispanser
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *handler) JoinRoom(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		return pkgErrors.WithStack(err)
	}

	roomName := c.Param("name")

	userId := c.Get("userId").(int)

	room, err := h.hub.GetRoomByName(roomName)
	if err != nil {
		room = &Room{
			Name:           roomName,
			Clients:        make(Clients),
			ColorDispanser: colordispanser.NewDispanser(),
		}

		h.hub.Rooms[roomName] = room
	}

	for _, client := range room.Clients {
		if client.Id == userId {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Fail to join in the room",
			})
		}
	}

	if len(room.Clients) >= 6 {
		return c.JSON(http.StatusForbidden, map[string]string{
			"message": "The room is full",
		})
	}

	user, err := h.storage.User.GetByID(userId)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	color, err := h.hub.Rooms[roomName].ColorDispanser.Get()
	if err != nil {
		return err
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		Id:       userId,
		RoomName: roomName,
		Username: user.Login,
		Color:    color,
	}

	m := &Message{
		Type:     "event",
		UserId:   cl.Id,
		RoomName: roomName,
		Content:  "A new user has joined the room",
		Color:    cl.Color,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub)

	return nil
}

type ClientResp struct {
	Color    string `json:"color"`
	Username string `json:"username"`
}

func (h *handler) GetClients(c echo.Context) error {
	clients := make([]ClientResp, 0)
	roomName := c.Param("name")

	if _, ok := h.hub.Rooms[roomName]; !ok {
		clients = make([]ClientResp, 0)
		return c.JSON(http.StatusOK, clients)
	}

	for _, client := range h.hub.Rooms[roomName].Clients {
		clients = append(clients, ClientResp{
			Color:    client.Color,
			Username: client.Username,
		})
	}

	return c.JSON(http.StatusOK, clients)
}
