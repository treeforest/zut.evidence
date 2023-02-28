package http

import (
	"github.com/gin-gonic/gin"
	"github.com/treeforest/zut.evidence/internal/service/file"
	"github.com/treeforest/zut.evidence/pkg/jwt"
)

type Server struct {
	engine *gin.Engine
	srv    *file.File
}

func New(addr string, f *file.File, jwtMgr *jwt.JWTManager) *Server {
	engine := gin.New()
	engine.Use(loggerHandler, recoverHandler, authHandler(jwtMgr))
	go func() {
		if err := engine.Run(addr); err != nil {
			panic(err)
		}
	}()
	s := &Server{
		engine: engine,
		srv:    f,
	}
	s.initRouter()
	return s
}

func (s *Server) initRouter() {
	v1 := s.engine.Group("/v1")

	fileGroup := v1.Group("/file")
	fileGroup.POST("/upload", s.fileUpload)
	fileGroup.GET("/download", s.fileDownload)
}

func (s *Server) Close() error {
	_ = s.srv.Close()
	return nil
}
