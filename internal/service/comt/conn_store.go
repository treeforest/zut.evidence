package comt

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"sync"
)

// connectionStore websocket 连接管理
type connectionStore struct {
	sync.RWMutex
	isClosing    bool
	wsId2Conn    map[int64]*connection
	shutdownOnce sync.Once
}

func newConnStore() *connectionStore {
	return &connectionStore{
		isClosing: false,
		wsId2Conn: make(map[int64]*connection),
	}
}

// getConnection 通过 websocket id 获取一个连接
func (cs *connectionStore) getConnection(id int64) (*connection, error) {
	cs.RLock()
	isClosing := cs.isClosing
	cs.RUnlock()

	if isClosing {
		return nil, errors.New("shutting down")
	}

	cs.RLock()
	conn, exists := cs.wsId2Conn[id]
	if exists {
		cs.RUnlock()
		return conn, nil
	}
	cs.RUnlock()

	return nil, errors.New("not found")
}

// shutdown 关闭所有连接
func (cs *connectionStore) shutdown() {
	cs.shutdownOnce.Do(func() {
		cs.Lock()
		cs.isClosing = true

		for _, conn := range cs.wsId2Conn {
			conn.close()
		}
		cs.wsId2Conn = make(map[int64]*connection)

		cs.Unlock()
	})
}

// closeById 关闭 websocket id 对应的连接
func (cs *connectionStore) closeById(id int64) {
	cs.Lock()
	defer cs.Unlock()
	if conn, exists := cs.wsId2Conn[id]; exists {
		conn.close()
		delete(cs.wsId2Conn, conn.id)
	}
}

// onConnected 存储连接
func (cs *connectionStore) onConnected(id int64, wsConn *websocket.Conn) *connection {
	cs.Lock()
	defer cs.Unlock()

	// 保持最新连接，若存在旧的连接，则关闭旧的连接
	if c, exists := cs.wsId2Conn[id]; exists {
		c.close()
	}

	conn := newConnection(id, wsConn)
	cs.wsId2Conn[id] = conn

	return conn
}
