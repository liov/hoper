package service

import (
	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/sdk/luosimao"
	"github.com/liov/hoper/server/go/user/confdao"
)

func LuosimaoVerify(vCode string) error {
	if err := luosimao.Verify(confdao.Conf.Customize.LuosimaoVerifyURL,
		confdao.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return errorcode.InvalidArgument.Message(err.Error())
	}
	return nil
}
