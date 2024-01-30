package service

import (
	"github.com/hopeio/tiga/protobuf/errorcode"
	"github.com/hopeio/tiga/utils/sdk/luosimao"
	"github.com/liov/hoper/server/go/user/confdao"
)

func LuosimaoVerify(vCode string) error {
	if err := luosimao.Verify(confdao.Conf.Customize.LuosimaoVerifyURL,
		confdao.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return errorcode.InvalidArgument.Message(err.Error())
	}
	return nil
}
