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
	ws, err := service.InitConnection(conn)
	if err != nil {
		log.Println(err)
	}
	go func() {
		for {
			data, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			ws.WriteMessage(data)
			fmt.Println("从客服端读取的消息为：" + string(data))
		}
	}()

	go func() {
		for {
			err := ws.WriteMessage([]byte("heart beat!!!!"))
			if err != nil {
				log.Println(err)
				return
			}
			time.Sleep(time.Second)
		}
	}()
}
