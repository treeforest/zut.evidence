package gateway

import (
	log "github.com/treeforest/logger"
	"net/http"
	"time"
)

func Logger(w http.ResponseWriter, r *http.Request, next func()) {
	// Start timer
	start := time.Now()
	path := r.URL.Path
	raw := r.URL.RawQuery
	method := r.Method
	addr := r.RemoteAddr

	// handle next
	next()

	// Stop timer
	end := time.Now()
	latency := end.Sub(start)
	if raw != "" {
		path = path + "?" + raw
	}
	log.Infof("METHOD:%s | PATH:%s | IP:%s | TIME:%d", method, path, addr, latency/time.Millisecond)
}

func Cors(w http.ResponseWriter, r *http.Request) bool {
	method := r.Method
	origin := r.Header.Get("Origin") //请求头部
	if origin != "" {
		//接收客户端发送的origin （重要！）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//服务器支持的所有跨域请求的方法
		w.Header().Set("Access-Control-Allow-Methods", "*")
		//允许跨域设置可以返回其他子段，可以自定义字段
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		w.Header().Set("Access-Control-Expose-Headers", "*")
		//设置缓存时间
		w.Header().Set("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	//允许类型校验
	if method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok!"))
		return true
	}
	return false
}
