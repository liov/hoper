package service

import (
	"github.com/hopeio/protobuf/errcode"
	"github.com/hopeio/utils/sdk/luosimao"
	"github.com/liov/hoper/server/go/user/confdao"
)

func LuosimaoVerify(vCode string) error {
	if err := luosimao.Verify(confdao.Conf.Customize.LuosimaoVerifyURL,
		confdao.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return errcode.InvalidArgument.Message(err.Error())
	}
	return nil
}
