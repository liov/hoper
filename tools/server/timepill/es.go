package timepill

import "time"

const DiaryIndex = "diary"

type IndexDiary struct {
	Id              int       `json:"id"`
	UserId          int       `json:"user_id" gorm:"index"`
	NoteBookId      int       `json:"notebook_id" gorm:"index"`
	NoteBookSubject string    `json:"notebook_subject" gorm:"index"`
	Content         string    `json:"content" gorm:"type:text"`
	Created         time.Time `json:"created" gorm:"timestamptz(6);default:'0001-01-01 00:00:00';index"`
}

func (diary *Diary) DiaryIndex() *IndexDiary {
	created, _ := time.Parse("2006-01-02 15:04:05", diary.Created)
	return &IndexDiary{
		Id:              diary.Id,
		UserId:          diary.UserId,
		NoteBookId:      diary.NoteBookId,
		NoteBookSubject: diary.NoteBookSubject,
		Content:         diary.Content,
		Created:         created,
	}
}
