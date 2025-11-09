package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait        = 2 * time.Second
	pongWait         = 60 * time.Second
	closeGracePeriod = 10 * time.Second
)

func readPump(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		log.Println(string(p))
	}
}

func writePump(conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(3000 * time.Millisecond)
LOOP:
	for {
		select {
		case <-done:
			log.Println("writer close...")
			break LOOP
		case t := <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("message from client %v", t))); err != nil {
				log.Println(err)
				break LOOP
			}
		}
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	conn.Close()
}

func main() {
	dialer := &websocket.Dialer{}
	conn, resp, err := dialer.Dial("ws://server:18443/ws", nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Status)

	conn.WriteMessage(websocket.TextMessage, []byte("first message from client"))

	done := make(chan struct{})

	go readPump(conn, done)
	go writePump(conn, done)

	<-done
}
