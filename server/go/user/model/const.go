package model

const (
	DefaultAvatar = "/static/images/6cbeb5c8-7160-4b6f-a342-d96d3c00367a.jpg"
)

const (
	UserTableName          = `user`
	UserExtTableName       = "user_ext"
	ResumeTableName        = "resume"
	UserActionLogTableName = "user_action_log"
	FollowTableName        = "user_follow"
)

const (
	ActiveDuration           = 24 * 60 * 60
	ResetDuration            = 24 * 60 * 60
	VerificationCodeDuration = 5 * 60
)
