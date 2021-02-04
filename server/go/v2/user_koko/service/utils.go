package service

import (
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/utils/verification"
	"github.com/liov/hoper/go/v2/utils/verification/validator"
)

func LuosimaoVerify(vCode string) error{
	if err := verification.LuosimaoVerify(conf.Conf.Customize.LuosimaoVerifyURL,
		conf.Conf.Customize.LuosimaoAPIKey, vCode); err != nil {
		return  errorcode.InvalidArgument.Message(err.Error())
	}
	return nil
}


func Validate(req interface{}) error {
	if err := validator.Validator.Struct(req); err != nil {
		return errorcode.InvalidArgument.Message(validator.Trans(err))
	}
	return nil
}
