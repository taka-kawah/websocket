package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	w.Write([]byte("this is ws"))
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
