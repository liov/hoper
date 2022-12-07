package service

import (
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	"github.com/liov/hoper/server/go/lib/utils/verification"
	"github.com/liov/hoper/server/go/mod/user/conf"
)

func LuosimaoVerify(vCode string) error {
	if err := verification.LuosimaoVerify(conf.Conf.Customize.LuosimaoVerifyURL,
		conf.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return errorcode.InvalidArgument.Message(err.Error())
	}
	return nil
}
