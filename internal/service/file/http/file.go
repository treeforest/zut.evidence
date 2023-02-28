package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/treeforest/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) fileUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Warn(err)
		errors(c, RequestErr, err.Error())
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		log.Warn(err)
		errors(c, RequestErr, err.Error())
		return
	}
	defer file.Close()
	cid, err := s.srv.UploadFile(c.GetInt64("uid"), fileHeader.Filename, file)
	if err != nil {
		log.Warn(err)
		errors(c, ServerErr, err.Error())
		return
	}
	res := map[string]string{
		"cid": cid,
	}
	result(c, res, 0)
}

func (s *Server) fileDownload(c *gin.Context) {
	var arg struct {
		Cid string `form:"cid" binding:"required"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}

	rc, filename, err := s.srv.DownloadFile(c.GetInt64("uid"), arg.Cid)
	if err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	defer rc.Close()

	contentType := ""
	b := strings.Split(filename, ".")
	switch b[len(b)-1] {
	case "jpg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "img":
		contentType = "application/x-img"
	case "jpe", "jpeg":
		contentType = "image/jpeg"
	case "gif":
		contentType = "image/gif"
	case "txt":
		contentType = "text/plain"
	case "zip":
		contentType = "application/zip"
	case "pbf":
		contentType = "application/pdf"
	case "word":
		contentType = "application/msword"
	default:
		contentType = "application/octet-stream"
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(filename)))
	c.Header("Content-Type", contentType)
	c.Header("Filename", filename)
	if _, err = io.Copy(c.Writer, rc); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
}
