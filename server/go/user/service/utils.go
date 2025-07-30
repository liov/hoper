package service

import (
	"github.com/hopeio/scaffold/errcode"
	"github.com/hopeio/gox/sdk/luosimao"
	"github.com/liov/hoper/server/go/global"
)

func LuosimaoVerify(vCode string) error {
	if err := luosimao.Verify(global.Conf.User.LuosimaoVerifyURL,
		global.Conf.User.LuosimaoAPIKey, vCode); err != nil {
		return errcode.InvalidArgument.Msg(err.Error())
	}
	return nil
}
