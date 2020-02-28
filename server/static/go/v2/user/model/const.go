package modelconst

const (
	DefaultAvatar = "/static/images/6cbeb5c8-7160-4b6f-a342-d96d3c00367a.jpg"
)

const (
	// 用户名的最大长度
	MaxUserNameLen = 10

	// 用户名的最小长度
	MinUserNameLen = 3

	// 密码的最大长度
	MaxPassLen = 15

	// 密码的最小长度
	MinPassLen = 6

	// 个性签名最大长度
	MaxSignatureLen = 200

	// 居住地的最大长度
	MaxLocationLen = 200

	//个人简介的最大长度
	MaxIntroduceLen = 500
)

const (
	ActiveDuration           = 24 * 60 * 60
	ResetDuration            = 24 * 60 * 60
	VerificationCodeDuration = 5 * 60
)
