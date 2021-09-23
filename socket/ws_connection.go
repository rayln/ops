package socket

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type WsConnection struct {
	wsConnect              *websocket.Conn
	inChan                 chan []byte
	outChan                chan []byte
	closeChan              chan byte
	channelId              int        //唯一标识
	mutex                  sync.Mutex // 对closeChan关闭上锁
	isClosed               bool       // 防止closeChan被关闭多次
	Token 				   int  //token 唯一凭证
}

func InitConnection(wsConn *websocket.Conn) (conn *WsConnection, err error) {
	conn = &WsConnection{
		wsConnect:              wsConn,
		inChan:                 make(chan []byte, 1000),
		outChan:                make(chan []byte, 1000),
		closeChan:              make(chan byte, 1),
	}
	// 启动读协程
	go conn.ReadLoop()
	// 启动写协程
	go conn.WriteLoop()
	return
}

func InitConnectionOnly(wsConn *websocket.Conn) (conn *WsConnection, err error) {
	conn = &WsConnection{
		wsConnect:              wsConn,
		inChan:                 make(chan []byte, 1000),
		outChan:                make(chan []byte, 1000),
		closeChan:              make(chan byte, 1),
	}
	return
}

func (conn *WsConnection) IsClose() bool {
	return conn.isClosed
}

func (conn *WsConnection) ReadMessage() (data []byte, err error) {

	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *WsConnection) WriteMessage(data []byte) (err error) {

	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

//update
func (conn *WsConnection) WriteMessageType(message_type int, data []byte) (err error) {

	//select {
	//case conn.outChan <- data:
	//case <-conn.closeChan:
	//	err = errors.New("connection is closed")
	//}
	conn.wsConnect.WriteMessage(message_type, data)
	return
}

func (conn *WsConnection) Close() {
	// 线程安全，可多次调用
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}

}

// 内部实现
func (conn *WsConnection) ReadLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		//阻塞在这里，等待inChan有空闲位置。。。。
		select {
		case conn.inChan <- data:
		case <-conn.closeChan: // closeChan 感知 conn断开
			goto ERR
		}

	}

ERR:
	conn.Close()
}

func (conn *WsConnection) WriteLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConnect.WriteMessage(websocket.BinaryMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}

func (conn *WsConnection) SetChannelId(channelId int) {
	conn.channelId = channelId
}
func (conn *WsConnection) GetChannelId() int {
	return conn.channelId
}
