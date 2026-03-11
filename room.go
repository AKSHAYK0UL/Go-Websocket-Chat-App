package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	data   []byte
	sender *Client
}

type Room struct {
	userList  map[*Client]bool
	joinRoom  chan *Client
	leaveRoom chan *Client
	broadCast chan Message
}

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

func (rm *RoomManager) getOrCreateRoom(roomName string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, exists := rm.rooms[roomName]; exists {
		return room
	}
	room := newRoom()
	rm.rooms[roomName] = room
	return room

}

func newRoom() *Room {
	r := &Room{
		userList:  make(map[*Client]bool),
		joinRoom:  make(chan *Client),
		leaveRoom: make(chan *Client),
		broadCast: make(chan Message),
	}
	go r.run()
	return r

}

func (r *Room) run() {
	for {
		select {
		case newUser := <-r.joinRoom:
			r.userList[newUser] = true
			go func() {
				r.broadCast <- Message{
					data:   []byte(newUser.userName + " has joined the room"),
					sender: newUser,
				}
			}()
		case leaveUser := <-r.leaveRoom:
			delete(r.userList, leaveUser)
			close(leaveUser.receiveMsg)
			go func() {
				r.broadCast <- Message{
					data:   []byte(leaveUser.userName + " has left the room"),
					sender: leaveUser,
				}

			}()
		case msg := <-r.broadCast:
			for c := range r.userList {
				if c != msg.sender {
					c.receiveMsg <- msg.data
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (rm *RoomManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("room")
	userName := r.URL.Query().Get("username")
	if roomName == "" || userName == "" {
		http.Error(w, "room and username are required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	room := rm.getOrCreateRoom(roomName)
	user := &Client{
		socket:     conn,
		userName:   userName,
		receiveMsg: make(chan []byte, 1024),
		room:       room,
	}
	room.joinRoom <- user
	defer func() {
		room.leaveRoom <- user

	}()
	go user.receive()
	user.send()
}
