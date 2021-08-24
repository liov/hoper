package main

import (
	"github.com/liov/hoper/server/go/mod/protobuf/content"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
	"github.com/liov/hoper/server/go/mod/tools/create-table/get"
	"github.com/liov/hoper/server/go/mod/upload/model"
)

var userMod = []interface{}{
	&user.User{},
	&user.UserExt{},
	&user.UserActionLog{},
	&user.UserBannedLog{},
	&user.UserFollow{},
	&user.UserScoreLog{},
	&user.UserDeviceInfo{},
	&user.Resume{},
}

var contentMod = []interface{}{
	/*	&content.Moment{},
		&content.Category{},
		&content.Tag{},
		&content.Comment{},
		&content.Like{},
		&content.ContentDel{},
		&content.Favorites{},
		&content.ContentExt{},

		&content.Report{},*/

	&content.Collection{},
}

var uploadMod = []interface{}{
	&model.UploadExt{},
	&model.UploadInfo{},
}

func main() {
	//get.GetDB().Debug().Migrator().DropTable(userMod...)
	//get.GetDB().Debug().Migrator().CreateTable(userMod...)
	get.GetDB().Debug().Migrator().CreateTable(contentMod...)
	//get.GetDB().Debug().Migrator().CreateTable(uploadMod...)
	//get.GetDB().Debug().Table("moment_comment").Migrator().CreateTable(&content.Comment{})
	//get.GetDB().Debug().Migrator().CreateTable(&model.UploadExt{})
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
