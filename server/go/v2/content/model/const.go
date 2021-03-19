package model

import "github.com/liov/hoper/go/v2/protobuf/content"

const (
	LikeTableName       = "like"
	BrowserTableName    = "browser"
	CommentTableName    = "comment"
	CollectTableName    = "collect"
	ReportTableName     = "report"
	GiveTableName       = "give"
	ApproveTableName    = "approve"
	MomentTableName     = "moment"
	NoteTableName     = "note"
	DiaryTableName = "diary"
	DiaryBookTableName = "diary_book"
	FavoritesTableName = "favorites"
	TagTableName        = "tag"
	ContentExtTableName = "content_ext"
)

const (
	HotSortSet = "Hot_Sort_Set"
)


func ActionTableName(action content.ActionType) string {
	switch action {
	case content.ActionBrowse:
		return BrowserTableName
	case content.ActionLike, content.ActionUnlike:
		return LikeTableName
	case content.ActionComment:
		return CommentTableName
	case content.ActionCollect:
		return CollectTableName
	case content.ActionReport:
		return ReportTableName
	case content.ActionGiveAction:
		return GiveTableName
	case content.ActionApprove:
		return ApproveTableName
	}
	return ""
}

func ContentTableName(typ content.ContentType) string {
	switch typ {
	case content.ContentMoment:
		return MomentTableName
	case content.ContentNote:
		return NoteTableName
	case content.ContentDairy:
		return DiaryTableName
	case content.ContentDairyBook:
		return DiaryBookTableName
	case content.ContentFavorites:
		return FavoritesTableName
	case content.ContentCollection:
		return CollectTableName
	case content.ContentComment:
		return CommentTableName
	}
	return ""
}