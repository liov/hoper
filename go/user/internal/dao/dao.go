package dao

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/bluele/gcache"
	"github.com/etcd-io/bbolt"
	"github.com/garyburd/redigo/redis"
	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
)

var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// DB 数据库连接
	db   *gorm.DB
	bolt *bbolt.DB
	// RedisPool Redis连接池
	redis *redis.Pool
	// MongoDB 数据库连接
	mongo       *mgo.Database
	redisExpire int32
	cache       gcache.Cache
	mcExpire    int32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Close close the resource.
func (d *dao) Close() {
	if d.bolt != nil {
		d.bolt.Close()
	}
	d.redis.Close()
	d.db.Close()
	if d.mongo != nil {
		d.mongo.Session.Close()
	}
}

func SetDB(db *gorm.DB) {
	db.Callback().Create().Remove("gorm:update_time_stamp")
	db.Callback().Update().Remove("gorm:update_time_stamp")
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")
	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Create().Replace("gorm:save_before_associations", saveBeforeAssociationsCallback)
	//db.Callback().Create().Replace("gorm:save_after_associations", saveAfterAssociationsCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	Dao.db = db
}

// updateTimeStampForCreateCallback will set `CreatedAt`, `ModifiedAt` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {

	if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
		if createdAtField.IsBlank {
			createdAtField.Set(time.Now())
		}
	}

}

// updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
		if updatedAtField.IsBlank {
			updatedAtField.Set(time.Now())
		}

	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedAtField, hasDeletedAtField := scope.FieldByName("deleted_at")

		if !scope.Search.Unscoped && hasDeletedAtField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(time.Now()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

func saveAssociationCheck(scope *gorm.Scope, field *gorm.Field) (autoUpdate bool, autoCreate bool, saveReference bool, r *gorm.Relationship) {
	checkTruth := func(value interface{}) bool {
		if v, ok := value.(bool); ok && !v {
			return false
		}

		if v, ok := value.(string); ok {
			v = strings.ToLower(v)
			if v == "false" || v != "skip" {
				return false
			}
		}

		return true
	}

	if changeableField(scope, field) && !field.IsBlank && !field.IsIgnored {
		if r = field.Relationship; r != nil {
			autoUpdate, autoCreate, saveReference = true, true, true

			if value, ok := scope.Get("gorm:save_associations"); ok {
				autoUpdate = checkTruth(value)
				autoCreate = autoUpdate
			} else if value, ok := field.TagSettings["SAVE_ASSOCIATIONS"]; ok {
				autoUpdate = checkTruth(value)
				autoCreate = autoUpdate
			}

			if value, ok := scope.Get("gorm:association_autoupdate"); ok {
				autoUpdate = checkTruth(value)
			} else if value, ok := field.TagSettings["ASSOCIATION_AUTOUPDATE"]; ok {
				autoUpdate = checkTruth(value)
			}

			if value, ok := scope.Get("gorm:association_autocreate"); ok {
				autoCreate = checkTruth(value)
			} else if value, ok := field.TagSettings["ASSOCIATION_AUTOCREATE"]; ok {
				autoCreate = checkTruth(value)
			}

			if value, ok := scope.Get("gorm:association_save_reference"); ok {
				saveReference = checkTruth(value)
			} else if value, ok := field.TagSettings["ASSOCIATION_SAVE_REFERENCE"]; ok {
				saveReference = checkTruth(value)
			}
		}
	}

	return
}

func changeableField(scope *gorm.Scope, field *gorm.Field) bool {
	if selectAttrs := scope.SelectAttrs(); len(selectAttrs) > 0 {
		for _, attr := range selectAttrs {
			if field.Name == attr || field.DBName == attr {
				return true
			}
		}
		return false
	}

	for _, attr := range scope.OmitAttrs() {
		if field.Name == attr || field.DBName == attr {
			return false
		}
	}

	return true
}

func saveBeforeAssociationsCallback(scope *gorm.Scope) {
	for _, field := range scope.Fields() {
		_, autoCreate, saveReference, relationship := saveAssociationCheck(scope, field)

		if relationship != nil && relationship.Kind == "belongs_to" {
			fieldValue := field.Field.Addr().Interface()
			//newScope := scope.New(fieldValue)

			if autoCreate {
				//scope.Err(scope.NewDB().Save(fieldValue).Error)
				//留个坑
				scope.NewDB().Save(fieldValue)
			}

			if saveReference {
				if len(relationship.ForeignFieldNames) != 0 {
					// set value's foreign key
					for idx, fieldName := range relationship.ForeignFieldNames {
						associationForeignName := relationship.AssociationForeignDBNames[idx]
						if foreignField, ok := scope.New(fieldValue).FieldByName(associationForeignName); ok {
							scope.Err(scope.SetColumn(fieldName, foreignField.Field.Interface()))
						}
					}
				}
			}
		}
	}
}

func saveAfterAssociationsCallback(scope *gorm.Scope) {
	for _, field := range scope.Fields() {
		_, autoCreate, saveReference, relationship := saveAssociationCheck(scope, field)

		if relationship != nil && (relationship.Kind == "has_one" || relationship.Kind == "has_many" || relationship.Kind == "many_to_many") {
			value := field.Field

			switch value.Kind() {
			case reflect.Slice:
				for i := 0; i < value.Len(); i++ {
					newDB := scope.NewDB()
					elem := value.Index(i).Addr().Interface()
					newScope := newDB.NewScope(elem)

					if saveReference {
						if relationship.JoinTableHandler == nil && len(relationship.ForeignFieldNames) != 0 {
							for idx, fieldName := range relationship.ForeignFieldNames {
								associationForeignName := relationship.AssociationForeignDBNames[idx]
								if f, ok := scope.FieldByName(associationForeignName); ok {
									scope.Err(newScope.SetColumn(fieldName, f.Field.Interface()))
								}
							}
						}

						if relationship.PolymorphicType != "" {
							scope.Err(newScope.SetColumn(relationship.PolymorphicType, relationship.PolymorphicValue))
						}
					}

					if autoCreate {

						//scope.Err(newDB.Save(elem).Error)
						newDB.Save(elem)

					}

					if !scope.New(newScope.Value).PrimaryKeyZero() && saveReference {
						if joinTableHandler := relationship.JoinTableHandler; joinTableHandler != nil {
							scope.Err(joinTableHandler.Add(joinTableHandler, newDB, scope.Value, newScope.Value))
						}
					}
				}
			default:
				elem := value.Addr().Interface()
				newScope := scope.New(elem)

				if saveReference {
					if len(relationship.ForeignFieldNames) != 0 {
						for idx, fieldName := range relationship.ForeignFieldNames {
							associationForeignName := relationship.AssociationForeignDBNames[idx]
							if f, ok := scope.FieldByName(associationForeignName); ok {
								scope.Err(newScope.SetColumn(fieldName, f.Field.Interface()))
							}
						}
					}

					if relationship.PolymorphicType != "" {
						scope.Err(newScope.SetColumn(relationship.PolymorphicType, relationship.PolymorphicValue))
					}
				}

				if newScope.PrimaryKeyZero() {
					if autoCreate {
						//scope.Err(scope.NewDB().Save(elem).Error)
						scope.NewDB().Save(elem)
					}
				}
			}
		}
	}
}

func SetCache(c gcache.Cache) {
	Dao.cache = c
}

func SetRedis(redb *redis.Pool) {
	Dao.redis = redb
}
