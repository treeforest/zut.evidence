package ipfs

import (
	shell "github.com/ipfs/go-ipfs-api"
	"io"
)

const (
	DefaultUrl = "localhost:5001"
)

type Ipfs struct {
	s *shell.Shell
}

func New(url string) *Ipfs {
	s := shell.NewShell(url)
	return &Ipfs{s: s}
}

func (o *Ipfs) Upload(r io.Reader) (cid string, err error) {
	return o.s.Add(r)
}

func (o *Ipfs) Download(cid string) (io.ReadCloser, error) {
	return o.s.Cat(cid)
}
