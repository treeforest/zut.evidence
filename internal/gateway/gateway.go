package gateway

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Gateway 网关
type Gateway struct {
	naming          *discovery.Discovery
	mux             *runtime.ServeMux
	handlers        []func(*runtime.ServeMux, *discovery.Discovery)
	hijackMap       map[string]http.HandlerFunc       // 记录劫持函数信息
	reverseProxyMap map[string]*httputil.ReverseProxy // 记录反向代理信息
}

func New(etcdUrl string) *Gateway {
	dis, err := discovery.New(etcdUrl)
	if err != nil {
		panic(err)
	}
	mux := runtime.NewServeMux()
	gw := &Gateway{
		naming:          dis,
		mux:             mux,
		handlers:        make([]func(*runtime.ServeMux, *discovery.Discovery), 0),
		hijackMap:       make(map[string]http.HandlerFunc),
		reverseProxyMap: make(map[string]*httputil.ReverseProxy),
	}
	return gw
}

func (gw *Gateway) Naming() *discovery.Discovery {
	return gw.naming
}

// Serve 启动网关服务
func (gw *Gateway) Serve(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	for _, handler := range gw.handlers {
		handler(gw.mux, gw.naming)
	}
	go func() {
		if err = http.Serve(ln, gw); err != nil {
			panic(err)
		}
	}()
	return nil
}

func (gw *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 异常恢复
	defer func() {
		if err := recover(); err != nil {
			log.Error("[Recover] ", err)
		}
	}()

	// 跨域配置
	if Cors(w, r) {
		return
	}

	handler := func() {
		// 劫持处理
		if hijack, ok := gw.hijackMap[hijackKey(r.Method, r.URL.Path)]; ok {
			log.Debug("hijack ", r.URL.Path)
			hijack(w, r)
			return
		}

		// 反向代理
		for prefix, proxy := range gw.reverseProxyMap {
			if strings.HasPrefix(r.URL.Path, prefix) {
				log.Debug("reverse proxy ", r.URL.Path)
				proxy.ServeHTTP(w, r)
				return
			}
		}

		// 正常处理
		log.Debug("normal ", r.URL.Path)
		gw.mux.ServeHTTP(w, r)
	}

	Logger(w, r, handler)
}

// Close 释放资源
func (gw *Gateway) Close() {
	gw.naming.Close()
}

// HandleFunc 注册处理函数
func (gw *Gateway) HandleFunc(handler func(*runtime.ServeMux, *discovery.Discovery)) {
	gw.handlers = append(gw.handlers, handler)
}

// HijackFunc 注册劫持函数，一些场景需要抛弃原先注册的处理器，使用劫持函数优先处理
func (gw *Gateway) HijackFunc(method, path string, handler http.HandlerFunc) {
	if _, ok := gw.hijackMap[hijackKey(method, path)]; ok {
		panic(fmt.Errorf("path [%s] already exist", path))
	}
	gw.hijackMap[hijackKey(method, path)] = handler
}

func hijackKey(method string, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

// ReverseProxyFunc 注册反向代理信息。prefix 表示需要进行反向代理的path的前缀，
// 例如 predix = "/v1/file", 则遇到 "/v1/file/upload" 时， 会直接采用反向代理处理请求。
// target 表示反向代理的目标地址。
func (gw *Gateway) ReverseProxyFunc(prefix string, rawUrl string) {
	target, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}
	if _, ok := gw.reverseProxyMap[prefix]; ok {
		panic(fmt.Errorf("prefix [%s] already exist", prefix))
	}
	gw.reverseProxyMap[prefix] = httputil.NewSingleHostReverseProxy(target)
}
