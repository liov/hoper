package global

import (
	"time"

	"github.com/hopeio/utils/os/fs"
)

type serverConfig struct {
	Volume   fs.Dir
	SiteURL  string
	PassSalt string
	// 天数
	TokenMaxAge      time.Duration
	TokenSecret      string
	TokenSecretBytes []byte
	PageSize         int8

	LuosimaoSuperPW   string
	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSaveDir fs.Dir //二维码保存路径
	FontSaveDir   fs.Dir //字体保存路径

	Limit Limit
}

type Limit struct {
	// 用户名的最大长度
	MaxUserNameLen uint8

	// 用户名的最小长度
	MinUserNameLen uint8

	// 密码的最大长度
	MaxPassLen uint8

	// 密码的最小长度
	MinPassLen uint8

	// 个性签名最大长度
	MaxSignatureLen uint8

	// 居住地的最大长度
	MaxLocationLen uint8

	//个人简介的最大长度
	MaxIntroduceLen uint8
}
