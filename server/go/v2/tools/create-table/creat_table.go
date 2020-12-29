package main

import (
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/dao/db/get"
)

var userMod = []interface{}{
	&model.User{},
	&model.UserExtend{},
	&model.UserActionLog{},
	&model.UserBannedLog{},
	&model.UserFollow{},
	&model.UserScoreLog{},
	&model.UserFollowLog{},
	&model.Resume{},
}

func main() {
	get.GetDB().Debug().Migrator().DropTable(userMod...)
	get.GetDB().Debug().Migrator().CreateTable(userMod...)
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
