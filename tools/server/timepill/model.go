package timepill

type User struct {
	Id       int      `json:"-"`
	UserId   int      `json:"id" gorm:"uniqueIndex:idx_id_name,priority:1"`
	Name     string   `json:"name" gorm:"uniqueIndex:idx_id_name,priority:2"`
	Intro    string   `json:"intro"`
	Created  string   `json:"created" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00';index"`
	State    int      `json:"state" gorm:"int2;default:0"`
	IconUrl  string   `json:"iconUrl" gorm:"size:255;default:''"`
	CoverUrl string   `json:"coverUrl" gorm:"size:255;default:''"`
	Badges   []*Badge `json:"badges" gorm:"-"`
}

type Badge struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id" gorm:"index"`
	BadgeId int    `json:"badge_id" gorm:"index"`
	Created string `json:"created" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00';index"`
	Title   string `json:"title" gorm:"size:255;default:''"`
	IconUrl string `json:"iconUrl" gorm:"size:255;default:''"`
}

type Diary struct {
	Id              int       `json:"id"`
	UserId          int       `json:"user_id" gorm:"index"`
	NoteBookId      int       `json:"notebook_id" gorm:"index"`
	NoteBookSubject string    `json:"notebook_subject" gorm:"index"`
	Content         string    `json:"content" gorm:"type:text"`
	Created         string    `json:"created" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00';index"`
	Updated         string    `json:"updated" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00'"`
	Type            int       `json:"type" gorm:"int2;default:0"`
	CommentCount    int       `json:"comment_count" gorm:"default:0"`
	PhotoUrl        string    `json:"photoUrl" gorm:"size:255;default:''"`
	PhotoThumbUrl   string    `json:"photoThumbUrl" gorm:"-"`
	LikeCount       int       `json:"like_count" gorm:"default:0"`
	Liked           bool      `json:"-" gorm:"-"`
	User            *User     `json:"user,omitempty" gorm:"-"`
	NoteBook        *NoteBook `json:"notebook,omitempty" gorm:"-"`
}

type TinyDiary struct {
	UserId   int
	PhotoUrl string
	Created  string
}

type NoteBook struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id" gorm:"index"`
	Subject     string `json:"subject" gorm:"size:255;index"`
	Description string `json:"description" gorm:"index"`
	Created     string `json:"created" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00';index"`
	Updated     string `json:"updated" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00'"`
	Expired     string `json:"expired" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00'"`
	Privacy     int    `json:"privacy" gorm:"int2;default:0"`
	Cover       int    `json:"cover" gorm:"int2;default:0"`
	CoverUrl    string `json:"coverUrl" gorm:"size:255;default:''"`
	IsPublic    bool   `json:"isPublic" gorm:"-"`
}

type Comment struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id" gorm:"index"`
	RecipientId int    `json:"recipient_id" gorm:"index"`
	DairyId     int    `json:"dairy_id" gorm:"index"`
	Content     string `json:"content" gorm:"type:text"`
	Created     string `json:"created" gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00';index"`
	User        *User  `json:"User" gorm:"-"`
	Recipient   *User  `json:"recipient" gorm:"-"`
}
