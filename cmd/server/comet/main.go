package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/treeforest/zut.evidence/internal/service"
	"github.com/treeforest/zut.evidence/internal/service/comt"
	grpcService "github.com/treeforest/zut.evidence/internal/service/comt/grpc"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"github.com/treeforest/zut.evidence/pkg/graceful"
	"github.com/treeforest/zut.evidence/pkg/jwt"

	log "github.com/treeforest/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	etcdUrl     *string              // 连接etcd的url
	port        *int                 // 服务监听端口
	serviceName string               // 服务名
	serviceAddr string               // 服务监听地址，由port计算得出
	jwtManager  *jwt.JWTManager      // jwt管理器
	naming      *discovery.Discovery // 注册与发现对象

	wsPort *int // websocket 端口
)

type ROLE = int

const (
	USER    ROLE = 1
	SCHOOL  ROLE = 2
	COMPANY ROLE = 3
	ROOT    ROLE = 4
)

// init 初始化函数(ps: 不需要改动)
func init() {
	etcdUrl = flag.String("etcdUrl", "localhost:2379", "the etcd url")
	port = flag.Int("port", 10004, "the grpc server port")
	wsPort = flag.Int("wsPort", 8082, "the websocket server port")
	flag.Parse()

	go func() {
		pprofAddr := fmt.Sprintf("127.0.0.1:%d", *port+1000)
		log.Infof("start Pprof at %s", pprofAddr)
		_ = http.ListenAndServe(pprofAddr, nil)
	}()

	serviceAddr = fmt.Sprintf("127.0.0.1:%d", *port)
	jwtManager = jwt.New(time.Second * jwt.DefaultTokenExpiration)

	var err error
	naming, err = discovery.New(*etcdUrl)
	if err != nil {
		panic(err)
	}
}

// accessibleRoles 访问控制列表(ps: 需要改动，自己填写访问控制列表)
// key:接口路径；value:标识哪些用户可以访问，若为空，则标识所有人都可以访问。
func accessibleRoles() map[string][]int {
	return map[string][]int{}
}

// register 服务注册(ps: 不需要改动)
func register(naming *discovery.Discovery, addr string) {
	err := naming.Register(context.TODO(), serviceName, addr)
	if err != nil {
		panic(err)
	}
}

// serverOptions grpc服务参数(ps: 不需要改动)
func serverOptions() []grpc.ServerOption {
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	return []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Minute * 30,
			Time:              time.Minute * 10,
			Timeout:           time.Second * 20,
		}),
	}
}

// Run 运行服务，并监听退出事件(ps: 不需要改动)
func Run(srv io.Closer, gSrv *grpc.Server) {
	register(naming, serviceAddr)
	graceful.Stop(func() {
		_ = naming.Close()
		gSrv.GracefulStop()
		if err := srv.Close(); err != nil {
			log.Errorf("server close error: %v", err)
		}
	})
}

func main() {
	// 服务初始化信息(ps: 需要改动)
	serviceName = "Comet"

	wsAddress := fmt.Sprintf("0.0.0.0:%d", *wsPort)
	srv := comt.InitWsServer(wsAddress, jwtManager)
	log.Infof("Start WS server at %s", wsAddress)

	grpcServer := grpcService.New(serviceAddr, srv, serverOptions()...)

	// 启动服务(ps: 不需要改动)
	log.Infof("Start GRPC server [%s] at %s", serviceName, serviceAddr)
	Run(srv, grpcServer)
}
