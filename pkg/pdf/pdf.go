package pdf

import (
	"bytes"
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/treeforest/zut.evidence/pkg/did"
	"image"
	"path"
)

// UtilPath 工具文件路径
var UtilPath = path.Join("pkg", "pdf")

// PdfInfo pdf中的内容信息
type PdfInfo struct {
	ProofClaim     *did.ProofClaim
	DownloadUrl    string // 履历下载链接
	ThumbImageData []byte // 缩略图
	AuthUrl        string // 认证地址链接
}

// GenEvidencePdf 生成存证PDF
func GenEvidencePdf(info *PdfInfo) []byte {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	W := gopdf.PageSizeA4.W
	H := gopdf.PageSizeA4.H

	// 添加字体-黑体
	err := pdf.AddTTFFont("simhei", path.Join(UtilPath, "simhei.ttf"))
	if err != nil {
		panic(err)
	}

	// 1 绘画边界
	drawBorder(pdf, W, H)

	// 2 标题
	err = pdf.SetFont("simhei", "", 24)
	if err != nil {
		panic(err)
	}
	drawTitle(pdf, W, 100)

	// 3 存证hash
	err = pdf.SetFont("simhei", "", 14)
	if err != nil {
		panic(err)
	}
	drawHash(pdf, W, 140, info.ProofClaim.Id)

	// 4 正文内容
	err = pdf.SetFont("simhei", "", 12)
	if err != nil {
		panic(err)
	}
	var leftX float64 = 160

	drawText(pdf, leftX, 280, fmt.Sprintf("证书作者: %s", info.ProofClaim.CredentialSubject.Id))
	drawText(pdf, leftX, 300, fmt.Sprintf("证书类型: %s", info.ProofClaim.CredentialSubject.Type))
	drawText(pdf, leftX, 320, fmt.Sprintf("证书简介: %s", info.ProofClaim.CredentialSubject.ShortDescription))
	drawText(pdf, leftX, 340, fmt.Sprintf("发行时间: %s", info.ProofClaim.IssuanceDate))
	drawText(pdf, leftX, 360, fmt.Sprintf("失效时间: %s", info.ProofClaim.ExpirationDate))
	drawText(pdf, leftX, 380, fmt.Sprintf("签名算法: %s", info.ProofClaim.Proof.Type))
	drawText(pdf, leftX, 400, fmt.Sprintf("发行人: %s", info.ProofClaim.Proof.Creator))
	// drawLink(pdf, leftX, 420, "履历原件", "点击下载", info.DownloadUrl)
	// 4.7 履历简略图
	//_ = drawThumb(pdf, W, leftX, 380, info.ThumbImageData)

	// 5 证书说明
	// 5.1 分割线
	pdf.SetLineType("dotted")
	pdf.Line(80, H-200, W-80, H-200)
	pdf.SetTextColor(0, 0, 0)
	// 5.2 证书说明
	drawText(pdf, 80, H-190, "证书说明: ")
	// 5.3 相关说明内容
	_ = pdf.SetFont("simhei", "", 10)
	drawText(pdf, 104, H-170, "本证书数据由高校联盟履历存证平台提供。")
	drawText(pdf, 104, H-155, "本证书可作为电子数据备案凭证。")
	drawText(pdf, 104, H-140, "如需验证证书的合法性与完整性，可通过认证机构进行查询和认证。")
	drawText(pdf, 104, H-125, "该证书已经在联盟链 FISCO BCOS 上进行存储。")

	// 6 电子印章
	pdf.Rotate(20, 300, 200)
	err = pdf.Image(path.Join(UtilPath, "seal.png"), 350, 200, nil)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	if err = pdf.Write(buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func drawBorder(pdf *gopdf.GoPdf, W, H float64) {
	// 设置边框
	pdf.SetStrokeColor(32, 58, 95)

	pdf.SetLineWidth(10)
	pdf.Line(0, 5, W, 5)
	pdf.Line(5, 0, 5, H)
	pdf.Line(0, H-5, W, H-5)
	pdf.Line(W-5, 0, W-5, H)

	pdf.SetLineWidth(2)
	pdf.Line(25, 26, W-25, 26)     // 上
	pdf.Line(26, 25, 26, H-25)     // 左
	pdf.Line(W-26, 25, W-26, H-25) // 右
	pdf.Line(25, H-26, W-25, H-26) // 下

	pdf.SetLineWidth(4)
	pdf.Line(29, 32, W-29, 32)     // 上
	pdf.Line(31, 30, 31, H-30)     // 左
	pdf.Line(W-31, 30, W-31, H-30) // 右
	pdf.Line(29, H-32, W-29, H-32) // 下

	pdf.SetLineWidth(1)
	pdf.Line(35, 37, W-35, 37)     // 上
	pdf.Line(35, 37, 35, H-37)     // 左
	pdf.Line(W-35, 37, W-35, H-37) // 右
	pdf.Line(35, H-37, W-35, H-37) // 下
}

func drawTitle(pdf *gopdf.GoPdf, W, y float64) {
	title := "区块链学业履历证书"
	titleWidth, _ := pdf.MeasureTextWidth(title)
	pdf.SetY(y)
	pdf.SetX((W - titleWidth) / 2)
	_ = pdf.Cell(nil, title)
}

func drawHash(pdf *gopdf.GoPdf, W, y float64, hash string) {
	hash = "证书ID: " + hash
	hashWidth, _ := pdf.MeasureTextWidth(hash)
	pdf.SetY(y)
	pdf.SetX((W - hashWidth) / 2)
	_ = pdf.Cell(nil, hash)
}

func drawText(pdf *gopdf.GoPdf, x, y float64, text string) {
	pdf.SetX(x)
	pdf.SetY(y)
	_ = pdf.Cell(nil, text)
}

func drawLink(pdf *gopdf.GoPdf, x, y float64, name, tip, url string) {
	name = name + ": "
	drawText(pdf, x, y, name)

	nameLen, _ := pdf.MeasureTextWidth(name)
	pdf.SetTextColor(0, 0, 204)
	drawText(pdf, x+nameLen, y, tip)

	tipLen, _ := pdf.MeasureTextWidth(tip)
	pdf.AddExternalLink(url, x+nameLen, y, tipLen, 12)

	pdf.SetTextColor(0, 0, 0)
}

func drawThumb(pdf *gopdf.GoPdf, W, x, y float64, imgData []byte) float64 {
	thumb := "履历缩略: "
	drawText(pdf, x, y, thumb)

	img, _, _ := image.Decode(bytes.NewReader(imgData))
	imgSize := img.Bounds().Size()

	thumbLen, _ := pdf.MeasureTextWidth(thumb)
	imgWidth := W - x*2 - thumbLen // 使得文字与图片的两边的间隔都是x
	imgHeight := (imgWidth / float64(imgSize.X)) * float64(imgSize.Y)
	x = x + thumbLen

	err := pdf.ImageFrom(img, x, y, &gopdf.Rect{W: imgWidth, H: imgHeight})
	if err != nil {
		panic(err)
	}

	// 返回图片的高度，便于后续内容的定位操作
	return imgHeight
}
