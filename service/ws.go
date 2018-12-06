package service

import (
	"github.com/gorilla/websocket"
	"log"
)

//---------------提供线程安全的WS服务-------------------
type WS struct {
	in    chan []byte // 入消息管道
	out   chan []byte // 出消息管道
	wsCon *websocket.Conn
}

func InitConnection(ws *websocket.Conn) (con *WS, err error) {
	con = &WS{
		wsCon: ws,
		in:    make(chan []byte, 100),
		out:   make(chan []byte, 100),
	}
	go readLoop(con)  // 拉起 读取消息的协程
	go writeLoop(con) // 拉起 发送消息的协程
	return
}

// API
func (ws *WS) ReadMessage() (data []byte, err error) {
	data = <-ws.in
	return
}

// API
func (ws *WS) WriteMessage(data []byte) (err error) {
	ws.out <- data
	return
}

// TODO:// 关闭ws链接
func (ws *WS) Close() {
	// 可重入的 ，线程安全
	ws.Close()
}

func readLoop(ws *WS) {
	for {
		_, data, err := ws.wsCon.ReadMessage()
		if err != nil {
			log.Fatal("读取ws数据失败")
			ws.Close()
		}
		// 从客户端读取的数据过多会阻塞在这里
		ws.in <- data
	}
}
func writeLoop(ws *WS) {
	for {
		data := <-ws.out
		err := ws.wsCon.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Fatal("读取ws数据失败")
			ws.Close()
		}
	}
}
