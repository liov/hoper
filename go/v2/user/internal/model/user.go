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
	UserSexMale = iota // 男

	UserSexFemale // 女

	UserSexNil //未填写
)

//用户操作
const(
	SignUp = iota
	ModifyPassword
)