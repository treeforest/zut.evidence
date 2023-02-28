package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/treeforest/zut.evidence/blockchain/contracts/evidence"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"path"
	"time"

	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
	"github.com/treeforest/zut.evidence/internal/service/logic"
	"github.com/treeforest/zut.evidence/internal/service/logic/conf"
	grcpService "github.com/treeforest/zut.evidence/internal/service/logic/grpc"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"github.com/treeforest/zut.evidence/pkg/graceful"
	"github.com/treeforest/zut.evidence/pkg/jwt"

	fiscobcosClient "github.com/FISCO-BCOS/go-sdk/client"
	fiscobcosConf "github.com/FISCO-BCOS/go-sdk/conf"
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
	port = flag.Int("port", 10005, "the grpc server port")
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
	requiredRoles := []int{USER, SCHOOL, COMPANY, ROOT}
	return map[string][]int{
		"/Logic/ApplyKYC":       requiredRoles,
		"/Logic/Issue":          requiredRoles,
		"/Logic/GetIssued":      requiredRoles,
		"/Logic/RevokeIssued":   requiredRoles,
		"/Logic/GetIssuerList":  requiredRoles,
		"/Logic/ChallengeSend":  requiredRoles,
		"/Logic/ChallengeReply": requiredRoles,
		"/Logic/ChallengeSent":  requiredRoles,
		"/Logic/ChallengeDoing": requiredRoles,
		"/Logic/ChallengeDone":  requiredRoles,
		"/Logic/Apply":          requiredRoles,
		"/Logic/ApplyDoing":     requiredRoles,
		"/Logic/ApplyDone":      requiredRoles,
		"/Logic/ApplyCount":     requiredRoles,
		"/Logic/Audit":          requiredRoles,
		"/Logic/AuditCount":     requiredRoles,
		"/Logic/AuditDoing":     requiredRoles,
		"/Logic/AuditFailed":    requiredRoles,
		"/Logic/AuditDone":      requiredRoles,
	}
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
	walletClientConn, err := naming.Dial("Wallet")
	if err != nil {
		panic(err)
	}
	didResolverConn, err := naming.Dial("DIDResolver")
	if err != nil {
		panic(err)
	}
	cometConn, err := naming.Dial("Comet")
	if err != nil {
		panic(err)
	}
	fileConn, err := naming.Dial("File")
	if err != nil {
		panic(err)
	}

	serviceName = "Logic"
	srv := logic.New(conf.Default(), pb.NewDIDResolverClient(didResolverConn), pb.NewWalletClient(walletClientConn),
		pb.NewCometClient(cometConn), pb.NewFileClient(fileConn), getEvidenceSession())
	grpcServer := grcpService.New(serviceAddr, srv, serverOptions()...)

	// 启动服务(ps: 不需要改动)
	log.Infof("Start GRPC server [%s] at %s", serviceName, serviceAddr)
	Run(srv, grpcServer)
}

func getEvidenceSession() func() (*evidence.EvidenceSession, error) {
	type ContractAddress struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	type ContractInfo struct {
		Addresses []ContractAddress `json:"addresses"`
	}

	configs, err := fiscobcosConf.ParseConfigFile(path.Join("cmd", "server", "logic", "config.toml"))
	if err != nil {
		log.Fatal("parse config file error: ", err)
	}
	config := &configs[0]

	// 加载合约地址
	b, err := ioutil.ReadFile(path.Join("blockchain", "deploy", "address.json"))
	if err != nil {
		log.Fatal(err)
	}
	contractInfo := ContractInfo{Addresses: make([]ContractAddress, 0)}
	if err = json.Unmarshal(b, &contractInfo); err != nil {
		log.Fatal(err)
	}

	// 初始化合约对象
	evidenceAddr := common.HexToAddress(contractInfo.Addresses[0].Address)

	return func() (*evidence.EvidenceSession, error) {
		ccCli, err := fiscobcosClient.Dial(config)
		if err != nil {
			log.Errorf("client dial failed: %v", err)
			return nil, err
		}

		instance, err := evidence.NewEvidence(evidenceAddr, ccCli)
		if err != nil {
			log.Errorf("cannot new confirm object: %v", err)
			return nil, err
		}
		callOpts := ccCli.GetCallOpts()
		callOpts.Context = context.Background()

		return &evidence.EvidenceSession{
			Contract:     instance,
			CallOpts:     *callOpts,
			TransactOpts: *ccCli.GetTransactOpts(),
		}, nil
	}
}
