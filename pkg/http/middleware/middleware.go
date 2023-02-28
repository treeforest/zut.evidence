package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/treeforest/zut.evidence/pkg/jwt"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
	"time"

	log "github.com/treeforest/logger"
)

// Logger 日志中间件
func Logger(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	method := c.Request.Method

	// Process request
	c.Next()

	// Stop timer
	end := time.Now()
	latency := end.Sub(start)
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	if raw != "" {
		path = path + "?" + raw
	}
	log.Infof("METHOD:%s | PATH:%s | CODE:%d | IP:%s | TIME:%d", method, path, statusCode, clientIP, latency/time.Millisecond)
}

// Recover panic场景下恢复的中间件
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			//log.Errorf("panic:%+v", err)
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			httprequest, _ := httputil.DumpRequest(c.Request, false)
			pnc := fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006-01-02 15:04:05"), string(httprequest), err, buf)
			fmt.Print(pnc)
			//log.Error(pnc)
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}

// Cors 跨域中间件
func Cors(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		//接收客户端发送的origin （重要！）
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "*")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "*")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "*")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	//允许类型校验
	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "ok!")
	}

	c.Next()
}

// Auth token 认证中间件
func Auth(jwtMgr *jwt.JWTManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.Contains(path, "login") || strings.Contains(path, "register") {
			// 拦截，进行token验证
			var authHeader struct {
				Token string `header:"token"`
			}
			if err := c.BindHeader(&authHeader); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    -400,
					"message": "without token in header",
				})
				c.Abort()
				return
			}
			uid, role, platform, extra, err := jwtMgr.Verify(authHeader.Token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    -400,
					"message": "invalid token",
				})
				c.Abort()
				return
			}
			// 设置上下文元数据
			c.Set("uid", uid)
			c.Set("role", role)
			c.Set("platform", platform)
			c.Set("extra", extra)
		}
		c.Next()
	}
}
