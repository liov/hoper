package service

import (
	"github.com/hopeio/scaffold/errcode"
	"github.com/hopeio/utils/sdk/luosimao"
	"github.com/liov/hoper/server/go/user/global"
)

func LuosimaoVerify(vCode string) error {
	if err := luosimao.Verify(global.Conf.Customize.LuosimaoVerifyURL,
		global.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return errcode.InvalidArgument.Msg(err.Error())
	}
	return nil
}
