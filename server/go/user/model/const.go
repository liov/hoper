package model

import "time"

const (
	DefaultAvatar = "image/6cbeb5c8-7160-4b6f-a342-d96d3c00367a.jpg"
)

const (
	UserTableName          = `user`
	UserExtTableName       = "user_ext"
	ResumeTableName        = "resume"
	UserActionLogTableName = "user_action_log"
	FollowTableName        = "user_follow"
)

const (
	ActiveDuration           = 24 * 60 * 60 * time.Second
	ResetDuration            = 24 * 60 * 60 * time.Second
	VerificationCodeDuration = 5 * 60 * time.Second
)
