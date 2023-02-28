package file

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/treeforest/zut.evidence/internal/service/file/conf"
	"github.com/treeforest/zut.evidence/internal/service/file/dao"
	"io"
)

// File 文件处理对象
type File struct {
	dao *dao.Dao
}

func New(c *conf.Config) *File {
	f := &File{
		dao: dao.New(c),
	}
	return f
}

func (f *File) Close() error {
	f.dao.Close()
	return nil
}

// UploadFile 上传文件
func (f *File) UploadFile(uid int64, filename string, r io.Reader) (string, error) {
	cid, err := f.dao.UploadFile(r)
	if err != nil {
		return "", errors.WithStack(err)
	}
	uid = 101 // for all
	if err = f.dao.AddFileMapping(uid, cid, filename); err != nil {
		return "", errors.WithStack(err)
	}
	return cid, nil
}

// DownloadFile 下载文件
func (f *File) DownloadFile(uid int64, cid string) (io.ReadCloser, string, error) {
	uid = 101 // for all
	filename, err := f.dao.ExistFileMapping(uid, cid)
	if err != nil {
		if err == redis.ErrNil {
			return nil, "", errors.New("文件不存在")
		}
		return nil, "", errors.WithStack(err)
	}
	rc, err := f.dao.DownloadFile(cid)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}
	return rc, filename, nil
}
