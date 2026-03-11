package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	socket     *websocket.Conn
	userName   string
	receiveMsg chan []byte
	room       *Room
}

func (c *Client) send() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}
		formatted := []byte(c.userName + ": " + string(msg))
		c.room.broadCast <- Message{data: formatted, sender: c}
	}
}

func (c *Client) receive() {
	defer c.socket.Close()

	for msg := range c.receiveMsg {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
