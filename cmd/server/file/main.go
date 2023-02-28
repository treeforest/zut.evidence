package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"io"
	pprofHttp "net/http"
	_ "net/http/pprof"
	"time"

	"github.com/treeforest/zut.evidence/internal/service/file"
	"github.com/treeforest/zut.evidence/internal/service/file/conf"
	grpcService "github.com/treeforest/zut.evidence/internal/service/file/grpc"
	"github.com/treeforest/zut.evidence/internal/service/file/http"
	"github.com/treeforest/zut.evidence/pkg/graceful"
	"github.com/treeforest/zut.evidence/pkg/jwt"

	log "github.com/treeforest/logger"
)

var (
	etcdUrl     *string              // 连接etcd的url
	rpcPort     *int                 // 服务监听端口
	httpPort    *int                 // 服务监听端口
	serviceName string               // 服务名
	serviceAddr string               // 服务监听地址，由port计算得出
	jwtManager  *jwt.JWTManager      // jwt管理器
	naming      *discovery.Discovery // 注册与发现对象
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
	httpPort = flag.Int("httpPort", 8081, "http server rpcPort")
	etcdUrl = flag.String("etcdUrl", "localhost:2379", "the etcd url")
	rpcPort = flag.Int("rpcPort", 10005, "the grpc server rpcPort")
	flag.Parse()

	go func() {
		pprofAddr := fmt.Sprintf("127.0.0.1:%d", *rpcPort+1000)
		log.Infof("start Pprof at %s", pprofAddr)
		_ = pprofHttp.ListenAndServe(pprofAddr, nil)
	}()

	serviceAddr = fmt.Sprintf("127.0.0.1:%d", *rpcPort)
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
	//interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	return []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		//grpc.UnaryInterceptor(interceptor.Unary()),
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
	serviceName = "File"

	f := file.New(conf.Default())
	grpcServer := grpcService.New(serviceAddr, f, serverOptions()...)
	log.Infof("Start GRPC server [%s] at %s", serviceName, serviceAddr)

	httpAddr := fmt.Sprintf("127.0.0.1:%d", *httpPort)
	srv := http.New(httpAddr, file.New(conf.Default()), jwtManager)
	defer srv.Close()
	log.Infof("Start HTTP server [%s] at %s", serviceName, httpAddr)

	Run(f, grpcServer)
}
