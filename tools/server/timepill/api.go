package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"net/http"
	"strconv"
)

const baseUrl = "https://open.timepill.net/api"
const v2Url = "https://v2.timepill.net/api"

type User struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Intro    string   `json:"intro"`
	Created  string   `json:"created"`
	State    int      `json:"state"`
	IconUrl  string   `json:"iconUrl"`
	CoverUrl string   `json:"coverUrl"`
	Badges   []string `json:"badges"`
}

func getSelfInfo() *User {
	var selfInfo User
	err := getV2("/users/my", nil, &selfInfo)
	if err != nil {
		log.Error(err)
	}
	return &selfInfo
}

func getUserInfo(id int) *User {
	var selfInfo User
	err := getV2("/users/"+strconv.Itoa(id), nil, &selfInfo)
	if err != nil {
		log.Error(err)
	}
	return &selfInfo
}

type TodayDiariesReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	FirstId  string `json:"first_id"`
}

type TodayDiaries struct {
	Count    int      `json:"count"`
	Page     string   `json:"page"`
	PageSize string   `json:"page_size"`
	Diaries  []*Diary `json:"diaries"`
}

type Diary struct {
	Id              int    `json:"id"`
	UserId          int    `json:"user_id" gorm:"index"`
	NoteBookId      int    `json:"notebook_id" gorm:"index"`
	NoteBookSubject string `json:"notebook_subject" gorm:"index"`
	Content         string `json:"content" gorm:"type:text"`
	Created         string `json:"created" gorm:"timestamptz(6);default:'0001-01-01 00:00:00';index"`
	Updated         string `json:"updated" gorm:"timestamptz(6);default:'0001-01-01 00:00:00'"`
	Type            int    `json:"type" gorm:"default:0"`
	CommentCount    int    `json:"comment_count" gorm:"default:0"`
	PhotoUrl        string `json:"photoUrl" gorm:"size:255;default:''"`
	PhotoThumbUrl   string `json:"photoThumbUrl" gorm:"-"`
	LikeCount       int    `json:"like_count" gorm:"default:0"`
	User            User   `json:"User" gorm:"-"`
}

func getTodayDiaries(page, pageSize int, firstId string) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1("/diaries/today", &TodayDiariesReq{page, pageSize, firstId}, &todayDiaries)
	if err != nil {
		log.Error(err)
	}
	return &todayDiaries
}

func getV1(api string, param, result interface{}) error {
	return call(http.MethodGet, baseUrl+api, param, result)
}
func postV1(api string, param, result interface{}) error {
	return call(http.MethodPost, baseUrl+api, param, result)
}

func getV2(api string, param, result interface{}) error {
	return call(http.MethodGet, v2Url+api, param, result)
}
func postV2(api string, param, result interface{}) error {
	return call(http.MethodPost, v2Url+api, param, result)
}

func callV1(method, api string, param, result interface{}) error {
	return call(method, baseUrl+api, param, result)
}

func callV2(method, api string, param, result interface{}) error {
	return call(method, v2Url+api, param, result)
}

func call(method, api string, param, result interface{}) error {
	return client.NewRequest(api, method, param).SetHeader("Authorization", token).SetLogger(nil).Do(result)
}
