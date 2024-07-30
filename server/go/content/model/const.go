package model

import "github.com/liov/hoper/server/go/protobuf/content"

const (
	Schema              = "content."
	TableNameLike       = Schema + "like"
	TableNameBrowser    = Schema + "browser"
	TableNameComment    = Schema + "comment"
	TableNameCollect    = Schema + "collection"
	TableNameContainer  = Schema + "container"
	TableNameReport     = Schema + "report"
	TableNameGive       = Schema + "give"
	TableNameApprove    = Schema + "approve"
	TableNameMoment     = Schema + "moment"
	TableNameNote       = Schema + "note"
	TableNameDiary      = Schema + "diary"
	TableNameArticle    = Schema + "article"
	TableNameDiaryBook  = Schema + "diary_book"
	TableNameFavorite   = Schema + "favorite"
	TableNameTag        = Schema + "tag"
	TableNameTagGroup   = Schema + "tag"
	TableNameStatistics = Schema + "statistics"
	TableNameContentTag = Schema + "content_tag"
)

const (
	HotSortSet = "Hot_Sort_Set"
)

func ActionTableName(action content.ActionType) string {
	switch action {
	case content.ActionBrowse:
		return TableNameBrowser
	case content.ActionLike, content.ActionUnlike:
		return TableNameLike
	case content.ActionComment:
		return TableNameComment
	case content.ActionCollect:
		return TableNameCollect
	case content.ActionReport:
		return TableNameReport
	case content.ActionGive:
		return TableNameGive
	case content.ActionApprove:
		return TableNameApprove
	}
	return ""
}

func ContentTableName(typ content.ContentType) string {
	switch typ {
	case content.ContentMoment:
		return TableNameMoment
	case content.ContentNote:
		return TableNameNote
	case content.ContentDairy:
		return TableNameDiary
	case content.ContentDairyBook:
		return TableNameDiaryBook
	case content.ContentFavorites:
		return TableNameFavorite
	case content.ContentCollection:
		return TableNameCollect
	case content.ContentComment:
		return TableNameComment
	}
	return ""
}
