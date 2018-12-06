package service

import (
	"github.com/gorilla/websocket"
	"sync"
	"fmt"
	"errors"
)

//---------------提供线程安全的WS服务-------------------
type WS struct {
	in    chan []byte // 入消息管道
	out   chan []byte // 出消息管道
	close chan byte
	wsCon *websocket.Conn
}

/**
TODO:// 初始化WS服务
 */
func InitConnection(ws *websocket.Conn) (con *WS, err error) {
	con = &WS{
		wsCon: ws,
		in:    make(chan []byte, 100),
		out:   make(chan []byte, 100),
		close: make(chan byte, 1),
	}

	go readLoop(con)  // 拉起 读取消息协程
	go writeLoop(con) // 拉起 发送消息协程
	return
}

// API
func (ws *WS) ReadMessage() (data []byte, err error) {
	select {
	case data = <-ws.in:
	case <-ws.close:
		ws.Close()
		err = errors.New("关闭连接")
	}

	return
}

// API
func (ws *WS) WriteMessage(data []byte) (err error) {
	select {
	case ws.out <- data:
	case <-ws.close:
		ws.Close()
		err = errors.New("关闭连接")
	}

	return
}

var once sync.Once
// TODO:// 关闭ws链接
func (ws *WS) Close() {
	// 可重入的 ，线程安全
	ws.wsCon.Close()
	fmt.Println("ws连接关闭")
	// 不可重入的，
	once.Do(func() {
		fmt.Println("信号管道关闭")
		close(ws.close)
	})

}

func readLoop(ws *WS) {
	for {
		_, data, err := ws.wsCon.ReadMessage()
		if err != nil {
			fmt.Println("读取ws数据失败")
			ws.Close()
			return
		}
		// 从客户端读取的数据过多会阻塞在这里
		select {
		case ws.in <- data:
		case <-ws.close:
			ws.Close()
			return
		}
	}
}
func writeLoop(ws *WS) {
	var data []byte
	for {
		select {
		case data = <-ws.out:
		case <-ws.close:
			ws.Close()
			return
		}
		err := ws.wsCon.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Println("写入ws数据失败")
			ws.Close()
			return
		}
	}
}
