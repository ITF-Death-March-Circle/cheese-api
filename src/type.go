package main

import "github.com/gorilla/websocket"

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}
type subscription struct {
	conn *connection
	room string
}
type ByteBroadCast struct {
	Message []byte
	Type    int
	Conn    *websocket.Conn
}
type message struct {
	data []byte
	room string
}

type hub struct {
	rooms      map[string]map[*connection]bool
	broadcast  chan message
	notify     chan message
	register   chan subscription
	unregister chan subscription
}

// ここからはシステム関連の型

type ResponseRoomId struct {
	Id   string `json:"id"`
	Link string `json:"link"`
}

type ErrorObject struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
type SocketRequest struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Name    string `json:"name"` //UserIdとする
	RoomId  string `json:"room_id"`
	UserId  string `json:"user_id"`
}

type MessageObject struct {
	// UserId  string `json:"user_id"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Action  string `json:"action"`
	// 多分これは必要無い
	// RoomId  string `json:"room_id"`
}

type VoteResponse struct {
	Action string `json:"action"`
	Value  int    `json:"value"`
	Text   string `json:"text"`
	Count  int    `json:"count"`
}

type InitObject struct {
	RoomId string `json:"room_id"`
	UserId string `json:"user_id"`
}

type LevelObject struct {
	Level  int    `json:"level"`
	Action string `json:"action"`
}
