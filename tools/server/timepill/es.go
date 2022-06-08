package timepill

const DiaryIndex = "diary"

const Mapping = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": {
        "properties": {
			"id": {
                "type": "long"
            },
            "user_id": {
                "type": "keyword",
                "fields": {
                    "raw": {
                        "type": "long"
                    }
                }
            },
            "notebook_id": {
                "type": "keyword",
                "fields": {
                    "raw": {
                        "type": "long"
                    }
                }
            },
            "notebook_subject": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart",
                "store": true,
                "fielddata": true
            },
            "content": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
            },
            "created": {
                "type": "date"
            }
        }
    }
}`

type IndexDiary struct {
	Id              int    `json:"id"`
	UserId          int    `json:"user_id" gorm:"index"`
	NoteBookId      int    `json:"notebook_id" gorm:"index"`
	NoteBookSubject string `json:"notebook_subject" gorm:"index"`
	Content         string `json:"content" gorm:"type:text"`
	Created         string `json:"created" gorm:"timestamptz(6);default:'0001-01-01 00:00:00';index"`
}

func (diary *Diary) DiaryIndex() *IndexDiary {

	return &IndexDiary{
		Id:              diary.Id,
		UserId:          diary.UserId,
		NoteBookId:      diary.NoteBookId,
		NoteBookSubject: diary.NoteBookSubject,
		Content:         diary.Content,
		Created:         diary.Created,
	}
}
