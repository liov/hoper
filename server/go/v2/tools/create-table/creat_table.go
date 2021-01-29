package main

import (
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/dao/db/get"
)

var userMod = []interface{}{
	&user.User{},
	&user.UserExt{},
	&user.UserActionLog{},
	&user.UserBannedLog{},
	&user.UserFollow{},
	&user.UserScoreLog{},
	&user.UserFollowLog{},
	&user.Resume{},
}

var contentMod = []interface{}{
	&content.Moment{},
	&content.Category{},
	&content.Tag{},
	&content.Comment{},
	&content.Like{},
	&content.ContentDel{},
	&content.Favorites{},
	&content.Collection{},
	&content.CountStatistics{},
	&content.Mood{},
	&content.Report{},
	&content.Collection{},
}

func main() {
	//get.GetDB().Debug().Migrator().DropTable(userMod...)
	//get.GetDB().Debug().Migrator().CreateTable(userMod...)
	get.GetDB().Debug().Migrator().CreateTable(contentMod...)
/*	get.GetDB().Exec(`CREATE OR REPLACE FUNCTION del_tabs(username IN VARCHAR) RETURNS void AS $$
		DECLARE
		statements CURSOR FOR
		SELECT tablename FROM pg_tables
		WHERE tableowner = username AND schemaname = 'public';
		BEGIN
		FOR stmt IN statements LOOP
		EXECUTE 'DROP TABLE ' || quote_ident(stmt.tablename) || ' CASCADE;';
		END LOOP;
		END;
		$$ LANGUAGE plpgsql`)*/
}

//清空所有表
//SELECT truncate_tables('postgres');
//删除所有表
//SELECT del_tabs('postgres');
