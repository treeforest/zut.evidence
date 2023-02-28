package gateway

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net/http"
	"net/url"
)

// HijackGenerateKey 拦截生成公私钥函数，主要是为了给前端返回压缩数据包，
// 因为 GRPC 网关无法返回文件（也许有，但是我还不知道）
func HijackGenerateKey(d *discovery.Discovery) http.HandlerFunc {
	cc, err := d.Dial("Wallet")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		if token == "" {
			errors(w, RequestErr, "not found token")
			return
		}

		// 请求钱包服务
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", token))
		reply, err := pb.NewWalletClient(cc).DownloadKey(ctx, &emptypb.Empty{})

		if err != nil {
			errors(w, RequestErr, err.Error())
			return
		}

		// 组装压缩包， 并返回压缩文件
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		var files = []struct {
			Name string
			Body []byte
		}{
			{"public_key.pem", reply.PublicKey},
			{"private_key.pem", reply.PrivateKey},
		}

		for _, file := range files {
			f, err := zipWriter.Create(file.Name)
			if err != nil {
				errors(w, ServerErr, err.Error())
				return
			}
			if _, err = f.Write(file.Body); err != nil {
				errors(w, ServerErr, err.Error())
				return
			}
		}

		if err = zipWriter.Close(); err != nil {
			errors(w, ServerErr, err.Error())
			return
		}

		// 返回数据
		filename := "public_private.zip"
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(filename)))
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Filename", filename)
		w.WriteHeader(http.StatusOK)
		if _, err = io.Copy(w, buf); err != nil {
			errors(w, ServerErr, err.Error())
			return
		}
	}
}
