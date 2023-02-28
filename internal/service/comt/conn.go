package comt

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

func newConnection(id int64, conn *websocket.Conn) *connection {
	err := conn.SetReadDeadline(time.Time{})
	if err != nil {
		panic(err)
	}
	conn.SetReadLimit(1024)
	conn.SetPingHandler(nil)
	conn.SetPongHandler(nil)
	conn.SetCloseHandler(nil)

	return &connection{
		id:       id,
		conn:     conn,
		outBuff:  make(chan *msgSending, 256),
		handler:  func([]byte) {},
		stopChan: make(chan struct{}, 1),
	}
}

type msgSending struct {
	msg   []byte
	onErr func(error)
}

// connection is a middleman between the ws connection and the hub.
type connection struct {
	id       int64            // 连接对应的uid
	conn     *websocket.Conn  // ws 连接
	outBuff  chan *msgSending // 发送通道
	handler  func([]byte)     // 消息处理回调
	stopChan chan struct{}
	stopOnce sync.Once
}

func (conn *connection) close() {
	conn.stopOnce.Do(func() {
		close(conn.stopChan)
	})
}

func (conn *connection) send(msg []byte, onErr func(error)) {
	m := &msgSending{
		msg:   msg,
		onErr: onErr,
	}
	select {
	case conn.outBuff <- m:
	case <-conn.stopChan:
	}
}

func (conn *connection) serviceConnection() error {
	errChan := make(chan error, 1)
	msgChan := make(chan []byte, 256)
	defer close(msgChan)

	go conn.readPump(errChan, msgChan)

	go conn.writePump(errChan)

	for {
		select {
		case <-conn.stopChan:
			return nil
		case err := <-errChan:
			return err
		case msg := <-msgChan:
			conn.handler(msg)
		}
	}
}

// 接收消息协程
func (conn *connection) readPump(errChan chan error, msgChan chan []byte) {
	defer func() {
		recover()
	}() // msgsCh might be closed

	for {
		select {
		case <-conn.stopChan:
			return
		default:
			_, message, err := conn.conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			select {
			case msgChan <- message:
			case <-conn.stopChan:
				return
			}
		}
	}
}

// 发送协程
func (conn *connection) writePump(errChan chan<- error) {
	const writeWait = time.Second * 20
	const pingPeriod = time.Second * 30

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		recover() // outBuff might be closed
		ticker.Stop()
	}()

	for {
		select {
		case m := <-conn.outBuff:
			_ = conn.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := conn.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				go m.onErr(err)
				errChan <- err
				return
			}

			if _, err = w.Write(m.msg); err != nil {
				go m.onErr(err)
				errChan <- err
				return
			}

			if err = w.Close(); err != nil {
				go m.onErr(err)
				errChan <- err
				return
			}

		case <-ticker.C:
			_ = conn.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				errChan <- err
				return
			}

		case <-conn.stopChan:
			return
		}
	}
}
