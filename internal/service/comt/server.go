package comt

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/pkg/jwt"
	"net/http"
	"time"
)

func InitWsServer(addr string, jwtMgr *jwt.JWTManager) *Server {
	s := newServer(jwtMgr)
	mux := http.NewServeMux()
	mux.Handle("/ws", s)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatal(err)
		}
	}()
	return s
}

type Server struct {
	jwtMgr    *jwt.JWTManager
	upgrader  *websocket.Upgrader
	connStore *connectionStore
}

func newServer(jwtMgr *jwt.JWTManager) *Server {
	return &Server{
		jwtMgr: jwtMgr,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		connStore: newConnStore(),
	}
}

func (s *Server) Close() error {
	s.connStore.shutdown()
	return nil
}

func (s *Server) PushMsg(uid int64, msgType int, msg string) error {
	conn, err := s.connStore.getConnection(uid)
	if err != nil {
		return errors.WithStack(err)
	}

	type pushMessage struct {
		Type int    `json:"type"`
		Data string `json:"data"`
	}
	m := &pushMessage{Type: msgType, Data: msg}
	b, _ := json.Marshal(m)

	conn.send(b, func(err error) {
		log.Errorf("Push message failed: %v", err)
		s.connStore.closeById(uid)
	})

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// websocket 请求升级
	wsConn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("Upgrade failed: %v", err)
		return
	}
	defer func() { _ = wsConn.Close() }()

	// 读取token信息
	_ = wsConn.SetReadDeadline(time.Now().Add(time.Second * 4))
	_, token, err := wsConn.ReadMessage()
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("token: %s", string(token))

	// 身份验证
	uid, _, _, _, err := s.jwtMgr.Verify(string(token))
	if err != nil {
		_ = wsConn.WriteMessage(websocket.CloseGoingAway, []byte(err.Error()))
		return
	}

	// 判断是否重复连接
	if _, err = s.connStore.getConnection(uid); err == nil {
		_ = wsConn.WriteMessage(websocket.CloseGoingAway, []byte("重复连接"))
		return
	}

	log.Infof("ONLINE | UID:%d", uid)
	conn := s.connStore.onConnected(uid, wsConn)

	defer func() {
		s.connStore.closeById(uid)
		log.Infof("OFFLINE | UID:%d", uid)
	}()

	// 推送测试
	//go func() {
	//	ticker := time.NewTicker(time.Second)
	//	defer ticker.Stop()
	//	for range ticker.C {
	//		if err := s.PushMsg(uid, 1, "这是推送测试数据"); err != nil {
	//			return
	//		}
	//	}
	//}()

	if err = conn.serviceConnection(); err != nil {
		if !websocket.IsUnexpectedCloseError(err) {
			log.Error(err)
		}
	}
}

//func (s *Server)
