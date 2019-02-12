package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// clientはチャットを行なっている1人のユーザを表します。
type client struct {
	// socketはこのクライアントのためのWebSocketです
	socket *websocket.Conn
	// sendはメッセージが送られるチャンネルです。
	send chan []byte
	// room はこのクライアントが参加しているチャットルームです。
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		fmt.Print(msg)
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
