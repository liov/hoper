package main

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"net/http"
	"strconv"
)

const baseUrl = "https://open.timepill.net/api"
const v2Url = "https://v2.timepill.net/api"

type User struct {
	Id       int      `json:"-"`
	UserId   int      `json:"id" gorm:"uniqueIndex:idx_id_name,priority:1"`
	Name     string   `json:"name" gorm:"uniqueIndex:idx_id_name,priority:2"`
	Intro    string   `json:"intro"`
	Created  string   `json:"created" gorm:"timestamptz(6);default:'0001-01-01 00:00:00';index"`
	State    int      `json:"state" gorm:"int2;default:0"`
	IconUrl  string   `json:"iconUrl" gorm:"size:255;default:''"`
	CoverUrl string   `json:"coverUrl" gorm:"size:255;default:''"`
	Badges   []*Badge `json:"badges" gorm:"-"`
	IsRecord bool     `json:"-"`
}

type Badge struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id" gorm:"index"`
	BadgeId int    `json:"badge_id" gorm:"index"`
	Created string `json:"created" gorm:"timestamptz(6);default:'0001-01-01 00:00:00';index"`
	Title   string `json:"title" gorm:"size:255;default:''"`
	IconUrl string `json:"iconUrl" gorm:"size:255;default:''"`
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

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type TodayDiariesReq struct {
	Page
	FirstId string `json:"first_id"`
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
	Type            int    `json:"type" gorm:"int2;default:0"`
	CommentCount    int    `json:"comment_count" gorm:"default:0"`
	PhotoUrl        string `json:"photoUrl" gorm:"size:255;default:''"`
	PhotoThumbUrl   string `json:"photoThumbUrl" gorm:"-"`
	LikeCount       int    `json:"like_count" gorm:"default:0"`
	Liked           bool   `json:"liked" gorm:"-"`
	User            *User  `json:"User" gorm:"-"`
}

func getTodayDiaries(page, pageSize int, firstId string) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1("/diaries/today", &TodayDiariesReq{Page{page, pageSize}, firstId}, &todayDiaries)
	if err != nil {
		log.Error(err)
	}
	return &todayDiaries
}

func getTodayTopicDiaries(page, pageSize int, firstId string) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1("/topic/diaries", &TodayDiariesReq{Page{page, pageSize}, firstId}, &todayDiaries)
	if err != nil {
		log.Error(err)
	}
	return &todayDiaries
}

func getFollowDiaries(page, pageSize int, firstId string) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1("/diaries/follow", &TodayDiariesReq{Page{page, pageSize}, firstId}, &todayDiaries)
	if err != nil {
		log.Error(err)
	}
	return &todayDiaries
}

type NotebookDiaries struct {
	Count    int      `json:"count"`
	Page     string   `json:"page"`
	PageSize string   `json:"page_size"`
	Items    []*Diary `json:"items"`
}

func getNotebookDiaries(id, page, pageSize int) *NotebookDiaries {
	var notebookDiaries NotebookDiaries
	err := getV1(fmt.Sprintf("/notebooks/%d/diaries", id), &Page{page, pageSize}, &notebookDiaries)
	if err != nil {
		log.Error(err)
	}
	return &notebookDiaries
}

func getUserTodayDiaries(userId int) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1(fmt.Sprintf("/users/%d/diaries", userId), nil, &todayDiaries)
	if err != nil {
		log.Error(err)
	}
	return &todayDiaries
}

type Comment struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id" gorm:"index"`
	RecipientId int    `json:"recipient_id"`
	DairyId     int    `json:"dairy_id"`
	Content     string `json:"content" gorm:"type:text"`
	Created     string `json:"created" gorm:"timestamptz(6);default:'0001-01-01 00:00:00';index"`
	User        *User  `json:"User" gorm:"-"`
	Recipient   *User  `json:"recipient" gorm:"-"`
}

func getDiaryComments(diaryId int) []*Comment {
	var comments []*Comment
	err := getV1(fmt.Sprintf("/diaries/%d/comments", diaryId), nil, &comments)
	if err != nil {
		log.Error(err)
	}
	return comments
}

type NoteBook struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id" gorm:"index"`
	Subject     string `json:"subject" gorm:"size:255;index"`
	Description string `json:"description" gorm:"index"`
	Created     string `json:"created" gorm:"timestamptz(6);default:'0001-01-01 00:00:00';index"`
	Updated     string `json:"updated" gorm:"timestamptz(6);default:'0001-01-01 00:00:00'"`
	Expired     string `json:"expired" gorm:"timestamptz(6);default:'0001-01-01 00:00:00'"`
	Privacy     int    `json:"privacy" gorm:"int2;default:0"`
	Cover       int    `json:"cover" gorm:"int2;default:0"`
	CoverUrl    string `json:"coverUrl" gorm:"size:255;default:''"`
	IsPublic    bool   `json:"isPublic" gorm:"-"`
}

func getUserNotebooks(userId int) []*NoteBook {
	var notebooks []*NoteBook
	err := getV1(fmt.Sprintf("/users/%d/notebooks", userId), nil, &notebooks)
	if err != nil {
		log.Error(err)
	}
	return notebooks
}

func getRelationUsers(page, pageSize int) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1("/relation", &Page{page, pageSize}, &todayDiaries)
	if err != nil {
		log.Error(err)
	}
	return &todayDiaries
}

func getRelationReverseUsers(page, pageSize int) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1("/relation/reverse", &Page{page, pageSize}, &todayDiaries)
	if err != nil {
		log.Error(err)
	}
	return &todayDiaries
}

func deleteDiary(diaryId int) *Response {
	var res Response
	err := call(http.MethodDelete, baseUrl+fmt.Sprintf("/diaries/%d", diaryId), nil, &res)
	if err != nil {
		log.Error(err)
	}
	return &res
}

func getRelation(userId int) *TodayDiaries {
	var todayDiaries TodayDiaries
	err := getV1(fmt.Sprintf("/relation/%d", userId), nil, &todayDiaries)
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

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
