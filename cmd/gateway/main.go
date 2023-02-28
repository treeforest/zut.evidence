package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/gateway"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"github.com/treeforest/zut.evidence/pkg/graceful"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/treeforest/logger"
)

func main() {
	etcdUrl := flag.String("etcdUrl", "localhost:2379", "the etcd url")
	port := flag.Int("port", 8080, "the gateway server port")
	flag.Parse()

	go func() {
		pprofAddr := fmt.Sprintf("127.0.0.1:%d", *port+1000)
		log.Infof("start Pprof at %s", pprofAddr)
		_ = http.ListenAndServe(pprofAddr, nil)
	}()

	gw := gateway.New(*etcdUrl)

	// 添加GRPC服务的反向代理
	gw.HandleFunc(func(mux *runtime.ServeMux, d *discovery.Discovery) {
		cc, err := d.Dial("Account")
		if err != nil {
			panic(err)
		}
		err = pb.RegisterAccountHandlerClient(context.TODO(), mux, pb.NewAccountClient(cc))
		if err != nil {
			panic(err)
		}
	})
	gw.HandleFunc(func(mux *runtime.ServeMux, d *discovery.Discovery) {
		cc, err := d.Dial("Wallet")
		if err != nil {
			panic(err)
		}
		err = pb.RegisterWalletHandlerClient(context.TODO(), mux, pb.NewWalletClient(cc))
		if err != nil {
			panic(err)
		}
	})
	gw.HandleFunc(func(mux *runtime.ServeMux, d *discovery.Discovery) {
		cc, err := d.Dial("DIDResolver")
		if err != nil {
			panic(err)
		}
		err = pb.RegisterDIDResolverHandlerClient(context.TODO(), mux, pb.NewDIDResolverClient(cc))
		if err != nil {
			panic(err)
		}
	})
	gw.HandleFunc(func(mux *runtime.ServeMux, d *discovery.Discovery) {
		cc, err := d.Dial("Logic")
		if err != nil {
			panic(err)
		}
		err = pb.RegisterLogicHandlerClient(context.TODO(), mux, pb.NewLogicClient(cc))
		if err != nil {
			panic(err)
		}
	})

	// 添加劫持函数
	gw.HijackFunc("GET", "/v1/wallet/downloadkey", gateway.HijackGenerateKey(gw.Naming()))

	// 添加反向代理
	gw.ReverseProxyFunc("/v1/file", "http://127.0.0.1:8081") // 文件服务
	gw.ReverseProxyFunc("/ws", "http://127.0.0.1:8082")      // websocket 服务（推送服务）

	// 启动网关服务
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	_ = gw.Serve(address)
	log.Infof("Start Gateway server at %s", address)

	// 关闭服务
	graceful.Stop(func() {
		gw.Close()
	})
}
