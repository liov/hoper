package service

import (
	"github.com/liov/hoper/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/v2/user/conf"
	"github.com/liov/hoper/v2/utils/verification"
)

func LuosimaoVerify(vCode string) error {
	if err := verification.LuosimaoVerify(conf.Conf.Customize.LuosimaoVerifyURL,
		conf.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return errorcode.InvalidArgument.Message(err.Error())
	}
	return nil
}
