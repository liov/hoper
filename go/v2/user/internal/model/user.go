package model

const(
	DefaultAvatar = "/static/images/6cbeb5c8-7160-4b6f-a342-d96d3c00367a.jpg"
)

//用户角色
const (
	UserRoleNormal = iota // 普通用户

	UserRoleAdmin // 管理员

	UserRoleSuperAdmin //超级管理员
)

//用户状态
const (
	UserStatusInActive = iota // 未激活

	UserStatusActivated // 已激活

	UserStatusFrozen //已冻结
)

//用户性别
const (
	UserSexNil = iota//未填写

	UserSexMale  // 男

	UserSexFemale // 女
)

//用户操作
const(
	SignUp = iota
	ModifyPassword
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
	// 生成激活账号的链接
	ActiveTime = "ActiveTime_"

	// 生成重置密码的链接
	ResetTime = "ResetTime_"

	// 用户信息
	LoginUser = "LoginUser_"

	// 用户每分钟最多能发表的文章数
	MomentMinuteLimit = "MomentMinuteLimit_"

	// 用户每天最多能发表的文章数
	MomentDayLimit = "MomentDayLimit_"

	// 用户每分钟最多能发表的文章数
	ArticleMinuteLimit = "ArticleMinuteLimit_"

	// 用户每天最多能发表的文章数
	ArticleDayLimit = "ArticleDayLimit_"

	// 用户每分钟最多能发表的评论数
	CommentMinuteLimit = "CommentMinuteLimit_"

	// 用户每天最多能发表的评论数
	CommentDayLimit = "CommentDayLimit_"
)

const (
	ActiveDuration = 24 * 60 * 60
	ResetDuration  = 24 * 60 * 60
)