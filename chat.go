package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"chat/service"
	"time"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("test"))
	})
	http.HandleFunc("/msg", handler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	ws, _ := service.InitConnection(conn)
	go func() {
		for {
			data, _ := ws.ReadMessage()
			ws.WriteMessage(data)
			fmt.Println("从客服端读取的消息为：" + string(data))
		}
	}()

	go func() {
		for {
			ws.WriteMessage([]byte("heart beat!!!!"))
			time.Sleep(time.Second)
		}
	}()
}
