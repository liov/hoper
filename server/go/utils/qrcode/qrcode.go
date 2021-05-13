package qrcode

import (
	"crypto/md5"
	"encoding/hex"
	"image/jpeg"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

/*func GetQrCodePath() string {
	return initialize.Config.Server.QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return initialize.Config.Server.RuntimeRootPath + initialize.Config.Server.QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	return initialize.Config.Server.PrefixUrl + "/" + GetQrCodePath() + name
}*/

func GetQrCodeFileName(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	_, err := os.Stat(src)

	return !os.IsNotExist(err)
}

func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	src := path + name
	_, err := os.Stat(src)
	if os.IsNotExist(err) == true {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := os.Open(path + name)
		if err != nil {
			return "", "", err
		}
		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}
