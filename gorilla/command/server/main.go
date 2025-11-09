package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait        = 2 * time.Second
	pongWait         = 60 * time.Second
	pingPeriod       = (pongWait * 9) / 10
	closeGracePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{}

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
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		log.Println(string(p))
	}
}

func writePump(conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(writeWait)
LOOP:
	for {
		select {
		case <-done:
			log.Println("write close...")
			break LOOP
		case t := <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("message from server %v", t))); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
					log.Println(err)
				}
				break LOOP
			}
		}
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	conn.Close()
}

func ping(conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping", err)
			}
		case <-done:
			return
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to upgrade"))
		return
	}

	done := make(chan struct{})

	go readPump(conn, done)
	go writePump(conn, done)
	go ping(conn, done)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("this is websocket server. please get /ws"))
	})
	http.HandleFunc("/ws", serveWs)
	fmt.Println("start listening at 18443...")
	err := http.ListenAndServe(":18443", nil)
	if err != nil {
		log.Fatal(err)
	}
}
