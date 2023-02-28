package comt

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"github.com/treeforest/zut.evidence/pkg/jwt"
	"io"
	"testing"
	"time"
)

func Test_Server(t *testing.T) {
	jwtMgr := jwt.New(time.Hour)

	// 启动推送服务
	wsServer := InitWsServer("localhost:18082", jwtMgr)
	defer wsServer.Close()
	time.Sleep(time.Millisecond * 500)

	// 生成token
	uid := int64(101)
	token, err := jwtMgr.Generate(uid, 1, "web", []byte{})
	require.NoError(t, err)

	// 连接websocket服务
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial("ws://localhost:18082/ws", map[string][]string{"authorization": {token}})
	require.NoError(t, err)

	defer conn.Close()

	// 等待接收推送的消息
	done := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			_, message, err := conn.ReadMessage()
			if err != nil {
				require.Equal(t, io.EOF, err)
			}
			t.Logf("[%d] received: %s\n", i, message)
		}
		done <- struct{}{}
	}()

	// 开始进行服务端推送
	for i := 0; i < 10; i++ {
		wsServer.PushMsg(uid, 0, "this is a push message")
	}

	<-done
}
