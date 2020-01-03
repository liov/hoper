package errorcode

type ErrCode uint32

const (
	SUCCESS       ErrCode = 0
	ERROR         ErrCode = 1
	InvalidParams ErrCode = 2
	SysError      ErrCode = 3

	LoginError   ErrCode = 1000 + iota //用户名或密码错误
	LoginTimeout                       //登录超时
	InActive                           //未激活账号
	NoAuthority

	ExistTag        ErrCode = 10001
	NotExistTag     ErrCode = 10002
	NotExistArticle ErrCode = 10003

	AuthCheckTokenFail ErrCode = 20001 + iota
	AuthCheckTokenTimeout
	AuthToken
	Auth

	// 保存图片失败
	UploadSaveImageFail ErrCode = 30001
	// 检查图片失败
	UploadCheckImageFail ErrCode = 30002
	// 校验图片错误，图片格式或大小有问题
	UploadCheckImageFormat ErrCode = 30003
	//尝试次数过多
	TimeTooMuch ErrCode = 40001

	//redis 5
	RedisErr ErrCode = 50001
)

func (e ErrCode) String() string {
	msg, ok := MsgFlags[e]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

func (e ErrCode) Error() string {
	return e.String()
}

var MsgFlags = map[ErrCode]string{
	SUCCESS:                "ok",
	ERROR:                  "fail",
	InvalidParams:          "请求参数错误",
	ExistTag:               "已存在该标签名称",
	NotExistTag:            "该标签不存在",
	NotExistArticle:        "该文章不存在",
	AuthCheckTokenFail:     "Token鉴权失败",
	AuthCheckTokenTimeout:  "Token已超时",
	AuthToken:              "Token生成失败",
	Auth:                   "Token错误",
	UploadSaveImageFail:    "保存图片失败",
	UploadCheckImageFail:   "检查图片失败",
	UploadCheckImageFormat: "校验图片错误，图片格式或大小有问题",
	LoginError:             "用户名或密码错误",
	LoginTimeout:           "登录超时",
	InActive:               "未激活账号",
}

const (
	Add = iota + 6000
	Sub
)
