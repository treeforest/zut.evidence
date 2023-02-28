package discovery

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/treeforest/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

var (
	leaseTTL = int64(6)       // 6s 服务租约超时时间，单位秒
	prefix   = "grpc/service" // 服务名前缀
)

const (
	BaseUrl = "http://localhost:2379"
)

type Discovery struct {
	logger log.Logger
	client *clientv3.Client
}

func New(etcdUrl string) (*Discovery, error) {
	cli, err := clientv3.NewFromURL(etcdUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Discovery{
		logger: log.NewStdLogger(log.WithPrefix("Discovery"), log.WithLevel(log.ERROR)),
		client: cli,
	}, nil
}

// Dial 连接grpc服务
func (d *Discovery) Dial(serviceName string) (*grpc.ClientConn, error) {
	rs, _ := resolver.NewBuilder(d.client)
	opts := []grpc.DialOption{
		grpc.WithResolvers(rs),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			// 若超过 10min 无消息来往（链接空闲），就发送一个ping消息检查链接是否还存在
			Time: time.Minute * 10,
			// 若在 20s 时间内未收到ping的响应消息，则主动断开链接
			Timeout: time.Second * 20,
			// 没有 active conn 时也发送 ping
			PermitWithoutStream: true,
		}),
	}
	target := fmt.Sprintf("etcd:///%s/%s", prefix, serviceName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	return grpc.DialContext(ctx, target, opts...)
}

func (d *Discovery) Register(ctx context.Context, serviceName, addr string) error {
	d.logger.Debugf("serviceName:%s | addr:%s | msg: try to register...", serviceName, addr)

	// 1. 创建一个租约
	lease := clientv3.NewLease(d.client)

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	leaseResp, err := lease.Grant(cancelCtx, leaseTTL)
	if err != nil {
		return errors.WithStack(err)
	}

	// 2.  续约（使租约保活）
	alive, err := lease.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	// 3. 创建一个服务的节点管理对象
	em, err := endpoints.NewManager(d.client, prefix)
	if err != nil {
		return errors.WithStack(err)
	}

	// 4. 注册服务节点信息到etcd
	cancelCtx, cancel = context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	err = em.AddEndpoint(context.TODO(), fmt.Sprintf("%s/%s/%s", prefix, serviceName, uuid.New().String()),
		endpoints.Endpoint{Addr: addr}, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return errors.WithStack(err)
	}

	d.logger.Debugf("serviceName:%s | addr:%s | msg: register success", serviceName, addr)

	unRegister := func() {
		cancelCtx, cancel = context.WithTimeout(ctx, time.Second*4)
		defer cancel()
		_ = em.DeleteEndpoint(ctx, serviceName)
		_ = lease.Close()
	}

	go d.keepRegister(ctx, alive, unRegister, serviceName, addr)

	return nil
}

func (d *Discovery) keepRegister(
	ctx context.Context,
	alive <-chan *clientv3.LeaseKeepAliveResponse,
	unRegister func(),
	serviceName string,
	addr string) {

	for {
		select {
		case resp := <-alive:
			if resp != nil {
				// keep alive success
				continue
			}

			d.logger.Warn("keep alive failed!")

			unRegister()

			for i := 0; i < 10; i++ {
				err := d.Register(ctx, serviceName, addr)

				if err != nil {
					d.logger.Errorf("Register failed | serviceName:%s | addr:%s | error:%v", serviceName, addr, err)
					time.Sleep(time.Second)
					continue
				}

				// 重连成功
				return
			}

			// 重连失败
			d.logger.Fatal("Register failed | serviceName:%s | addr:%s", serviceName, addr)

		case <-ctx.Done():
			unRegister()
		}
	}

}

func (d *Discovery) Close() error {
	return d.client.Close()
}
