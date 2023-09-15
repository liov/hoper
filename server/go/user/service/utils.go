package service

import (
	"github.com/hopeio/lemon/protobuf/errorcode"
	"github.com/hopeio/lemon/utils/sdk/luosimao"
	"github.com/liovx/hoper/server/go/user/confdao"
)

func LuosimaoVerify(vCode string) error {
	if err := luosimao.Verify(confdao.Conf.Customize.LuosimaoVerifyURL,
		confdao.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return errorcode.InvalidArgument.Message(err.Error())
	}
	return nil
}
