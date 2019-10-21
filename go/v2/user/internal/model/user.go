package model

import "time"

type User struct {
	ID              uint64      `gorm:"primary_key" json:"id"`
	Name            string      `gorm:"type:varchar(10);not null" json:"name"`
	ActivatedAt     *time.Time  `json:"-"` //激活时间
	Password        string      `gorm:"type:varchar(100)" json:"-"`
	Account         string      `gorm:"type:varchar(20);unique_index" json:"-"`
	Email           string      `gorm:"type:varchar(20);unique_index;not null" json:"email"`
	Phone           *string     `gorm:"type:varchar(20);unique_index" json:"phone"` //手机号
	Sex             string      `gorm:"type:varchar(1);not null" json:"sex"`
	Birthday        *time.Time  `json:"birthday"`
	Introduction    string      `gorm:"type:varchar(500)" json:"introduction"` //简介
	Score           uint64      `gorm:"default:0" json:"score"`                //积分
	Signature       string      `gorm:"type:varchar(100)" json:"signature"`    //个人签名
	Role            uint8       `gorm:"type:smallint;default:0" json:"-"`      //管理员or用户
	AvatarURL       string      `gorm:"type:varchar(100)" json:"avatar_url"`   //头像
	CoverURL        string      `gorm:"type:varchar(100)" json:"cover_url"`    //个人主页背景图片URL
	Address         string      `gorm:"type:varchar(100)" json:"address"`
	Location        string      `gorm:"type:varchar(100)" json:"location"`
	EduExps         []Education `json:"edu_exps"`  //教育经历
	WorkExps        []Work      `json:"work_exps"` //职业经历
	UpdatedAt       *time.Time  `json:"-"`
	BannedAt        *time.Time  `sql:"index" json:"banned_at"`
	CreatedAt       time.Time   `json:"-"`
	LastActivatedAt *time.Time  `json:"-"`                                  //上次活跃时间
	LastName        string      `gorm:"type:varchar(100)" json:"last_name"` //上个名字
	Status          uint8       `gorm:"type:smallint;default:0" json:"-"`   //状态
	Follows         []User      `gorm:"-" json:"follows"`                   //gorm:"foreignkey:FollowID []Follow里的User
	Followeds       []User      `gorm:"-" json:"followeds"`                 //gorm:"foreignkey:UserID"	[]Follow里的FollowUser
	FollowCount     uint64      `gorm:"default:0" json:"follow_count"`      //关注数量
	FollowedCount   uint64      `gorm:"default:0" json:"followed_count"`    //被关注数量
	ArticleCount    uint64      `gorm:"default:0" json:"article_count"`     //文章数量
	MomentCount     uint64      `gorm:"default:0" json:"moment_count"`
	DiaryBookCount  uint64      `gorm:"default:0" json:"diary_book_count"`
	DiaryCount      uint64      `gorm:"default:0" json:"diary_count"`
	CommentCount    uint64      `gorm:"default:0" json:"comment_count"` //评论数量
}

type Resume struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	Kind       uint8     `gorm:"type:smallint" json:"kind"`
	School     string    `gorm:"type:varchar(20)" json:"school"`
	Speciality string    `gorm:"type:varchar(100)" json:"speciality"` //专业
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	UserID     uint64    `json:"user_id"`
	Status     uint8     `gorm:"type:smallint;default:0" json:"-"`
}

type Education struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	School     string    `gorm:"type:varchar(20)" json:"school"`
	Speciality string    `gorm:"type:varchar(100)" json:"speciality"` //专业
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	UserID     uint64    `json:"user_id"`
	Status     uint8     `gorm:"type:smallint;default:0" json:"-"`
}

type Work struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	Company   string    `json:"company"` //公司或组织
	Title     string    `json:"title"`   //职位
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	UserID    uint64    `json:"user_id"`
	Status    uint8     `gorm:"type:smallint;default:0" json:"-"`
}
