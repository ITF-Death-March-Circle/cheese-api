package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//消せ
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s subscription) readPump() {
	c := s.conn
	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.EnableWriteCompression(true)
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	first_m, err := getCurrentValue()
	if err != nil {
		log.Fatalln(err)
	}
	h.broadcast <- message{first_m, s.room}
	for {
		_, msg, err := c.ws.ReadMessage()
		//ここでハンドラーを噛ませば各種処理を行える
		//Todo
		response := handler(msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := message{response, s.room}
		h.broadcast <- m
	}
}

//恐らく非同期にブロードキャストして良い

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}
}

// func sendStatusRoutines() {
// 	for {
// 		//ここでlevelをよしなに持ってくる
// 		level, _ := redis.GetValue("COMMON_VALUE_PATH")

// 		result := 1

// 		if level <= LEVEL_1 {
// 			result = 1
// 		} else if level <= LEVEL_2 {
// 			result = 2
// 		} else if level <= LEVEL_3 {
// 			result = 3
// 		} else if level <= LEVEL_4 {
// 			result = 4
// 		} else if level <= LEVEL_5 {
// 			result = 5
// 		} else {
// 			result = 6
// 		}
// 		// コールバックオブジェクトを作詞絵
// 		callBack := LevelObject{
// 			Action: NOTIFY_CURRENT_LEVEL,
// 			Level:  result,
// 		}
// 		b, err := json.Marshal(callBack)
// 		if err != nil {
// 			log.Println("cannot marshal struct: %v", err)
// 			return
// 		}
// 		m := message{b, "maid"}
// 		h.broadcast <- m
// 		//5秒待機
// 		time.Sleep(5 * time.Second)
// 	}
// }

func serveWs(w http.ResponseWriter, r *http.Request, roomId string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}
	h.register <- s
	go s.writePump()
	go s.readPump()
	// go sendStatusRoutines()
}
