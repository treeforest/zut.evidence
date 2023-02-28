package dao

import "io"

func (d *Dao) UploadFile(r io.Reader) (string, error) {
	return d.ipfs.Upload(r)
}

func (d *Dao) DownloadFile(cid string) (io.ReadCloser, error) {
	return d.ipfs.Download(cid)
}
