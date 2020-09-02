import 'dart:async';
import 'package:flutter_luakit_plugin/flutter_luakit_plugin.dart';

enum FieldType {
  Integer,
  Real,
  Blob,
  Char,
  Text,
  Boolean,
}

enum FilterType {
  WHERE_COLUMS,
  WHERE_SQL,
  PRIMARY_KEY,
  LIMIT,
  OFFSET,
  ORDER_BY,
  GROUP_BY,
  HAVING,
  HAVING_BY_BINDS,
  NEED_COLUMNS,
  JOIN,
}

enum WhereCondictionType {
  LESS_THEN,
  EQ_OR_LESS_THEN,
  MORE_THEN,
  EQ_OR_MORE_THEN,
  IN,
  NOT_IN,
  IS_NULL,
  LIKE,
}

enum JoinType {
  INNER ,
  LEFT ,
}




class FlutterOrmPlugin {

  static Future<void> createTable(String dbName, String tableName, Map<String,Field> fields) async {
    Map<String, dynamic> tableArgs = new Map<String, dynamic>();
    tableArgs["__dbname__"] = dbName;
    tableArgs["__tablename__"] = tableName;
    Map<String, dynamic> args = new Map<String, dynamic>();
    fields.forEach((key, value){
      Field f = value;
      tableArgs[key] = f.getParamsMap();
    });
    args["name"] = tableName;
    args["args"] = tableArgs;
    await FlutterLuakitPlugin.callLuaFun("orm.class.table", "addTableInfo", args);
  }

  static Future<dynamic> saveOrm(String tableName, Map value) async {
    Map<String, dynamic> args = new Map<String, dynamic>();
    args["name"] = tableName;
    args["args"] = value;
    return await FlutterLuakitPlugin.callLuaFun("orm.class.table", "saveOrm", args);
  }

  static Future<void> batchSaveOrms(String tableName, List values) async {
    Map<String, dynamic> args = new Map<String, dynamic>();
    args["name"] = tableName;
    args["args"] = values;
    await FlutterLuakitPlugin.callLuaFun("orm.class.table", "batchSaveOrms", args);
  }

}

class Filter {

  final FilterType type;

  final dynamic value;

  Filter(this.type,this.value);

}


class WhereCondiction{

  final String columName;

  final WhereCondictionType type;

  final dynamic value;

  WhereCondiction(this.columName, this.type,this.value);

}

class JoinCondiction{

  final String tableName;

  JoinType type = JoinType.INNER;

  String where;

  List<dynamic> whereBindingValues;

  List<String> needColumns;

  Map<String,String> matchColumns;

  JoinCondiction(this.tableName);

}



class Query{

  final String tableName;

  List<Filter> _filters = new List<Filter>();

  Query(this.tableName);

  Query whereByColumFilters(List<WhereCondiction> whereCondictions){
    _filters.add(Filter(FilterType.WHERE_COLUMS, whereCondictions));
    return this;
  }

  Query whereBySql(String sql, List<dynamic> args){
    Map<String,dynamic> m = new Map<String,dynamic>();
    m["sql"] = sql;
    if(args != null){
      m["args"] = args;
    }
    _filters.add(Filter(FilterType.WHERE_SQL, m));
    return this;
  }

  Query primaryKey(List<dynamic> values){
    _filters.add(Filter(FilterType.PRIMARY_KEY, values));
    return this;
  }

  Query limit(int value){
    _filters.add(Filter(FilterType.LIMIT, value));
    return this;
  }

  Query offset(int value){
    _filters.add(Filter(FilterType.OFFSET, value));
    return this;
  }

  Query clear(){
    _filters.clear();
    return this;
  }

  Query orderBy(List<String> orderStrings){
    _filters.add(Filter(FilterType.ORDER_BY, orderStrings));
    return this;
  }

  Query needColums(List<String> columNames){
    _filters.add(Filter(FilterType.NEED_COLUMNS, columNames));
    return this;
  }

  Query groupBy(List<String> groupByStrings){
    _filters.add(Filter(FilterType.GROUP_BY, groupByStrings));
    return this;
  }

  Query having(Map args){
    _filters.add(Filter(FilterType.HAVING, args));
    return this;
  }

  Query havingByBindings(String sql , List bindingValues){
    Map<String,dynamic> m = new Map<String,dynamic>();
    m["sql"] = sql;
    if(bindingValues != null){
      m["args"] = bindingValues;
    }
    _filters.add(Filter(FilterType.HAVING_BY_BINDS, m));
    return this;
  }

  Query join(JoinCondiction c){
    _filters.add(Filter(FilterType.JOIN, c));
    return this;
  }

  Future<List<dynamic>> all () async{
    Map<String, dynamic> args = new Map<String, dynamic>();
    args["name"] = tableName;
    args["args"] = _luakitSqlParams();
    List<dynamic> list = await FlutterLuakitPlugin.callLuaFun("orm.class.table", "getAllByParams", args);
    return list;
  }

  Future<Map> first() async{
    Map<String, dynamic> args = new Map<String, dynamic>();
    args["name"] = tableName;
    args["args"] = _luakitSqlParams();
    Map orm = await FlutterLuakitPlugin.callLuaFun("orm.class.table", "getFirstByParams", args);
    return orm;
  }

  Future<void> update(Map<String, dynamic> updateValue) async {
    Map<String, dynamic> args = new Map<String, dynamic>();
    args["name"] = tableName;
    args["args"] = _luakitSqlParams();
    args["updateValue"] = updateValue;
    await FlutterLuakitPlugin.callLuaFun("orm.class.table", "updateByParams", args);
  }

  Future<void> delete() async {
    Map<String, dynamic> args = new Map<String, dynamic>();
    args["name"] = tableName;
    args["args"] = _luakitSqlParams();
    await FlutterLuakitPlugin.callLuaFun("orm.class.table", "deleteByParams", args);
  }

  List<dynamic> _luakitSqlParams()
  {
    List<dynamic> params = new List<dynamic>();
    _filters.forEach((Filter f){
      Map<String,dynamic> m = new Map<String,dynamic>();
      if(f.type == FilterType.WHERE_COLUMS){
        m["type"] = "WHERE_COLUMS";
        List<WhereCondiction> l = f.value;
        List<Map<String,dynamic>> whereCondictions = new List<Map<String,dynamic>>();
        l.forEach((WhereCondiction c){
          Map<String,dynamic> whereCondiction = new Map<String,dynamic>();
          whereCondiction["columName"] = c.columName;
          whereCondiction["value"] = c.value;
          if(c.type == WhereCondictionType.LESS_THEN){
            whereCondiction["type"] = "LESS_THEN";
          } else if(c.type == WhereCondictionType.EQ_OR_LESS_THEN){
            whereCondiction["type"] = "EQ_OR_LESS_THEN";
          } else if(c.type == WhereCondictionType.MORE_THEN){
            whereCondiction["type"] = "MORE_THEN";
          } else if(c.type == WhereCondictionType.EQ_OR_MORE_THEN){
            whereCondiction["type"] = "EQ_OR_MORE_THEN";
          } else if(c.type == WhereCondictionType.IN){
            whereCondiction["type"] = "IN";
          } else if(c.type == WhereCondictionType.NOT_IN){
            whereCondiction["type"] = "NOT_IN";
          } else if(c.type == WhereCondictionType.IS_NULL){
            whereCondiction["type"] = "IS_NULL";
          } else if(c.type == WhereCondictionType.LIKE){
            whereCondiction["type"] = "LIKE";
          }
          whereCondictions.add(whereCondiction);
        });
        m["value"] = whereCondictions;
      } else if(f.type == FilterType.WHERE_SQL){
        m["type"] = "WHERE_SQL";
        m["value"] = f.value;
      } else if(f.type == FilterType.PRIMARY_KEY){
        m["type"] = "PRIMARY_KEY";
        m["value"] = f.value;
      } else if(f.type == FilterType.PRIMARY_KEY){
        m["type"] = "PRIMARY_KEY";
        m["value"] = f.value;
      } else if(f.type == FilterType.LIMIT){
        m["type"] = "LIMIT";
        m["value"] = f.value;
      } else if(f.type == FilterType.OFFSET){
        m["type"] = "OFFSET";
        m["value"] = f.value;
      } else if(f.type == FilterType.ORDER_BY){
        m["type"] = "ORDER_BY";
        m["value"] = f.value;
      } else if(f.type == FilterType.GROUP_BY){
        m["type"] = "GROUP_BY";
        m["value"] = f.value;
      } else if(f.type == FilterType.HAVING){
        m["type"] = "HAVING";
        m["value"] = f.value;
      } else if(f.type == FilterType.HAVING_BY_BINDS){
        m["type"] = "HAVING_BY_BINDS";
        m["value"] = f.value;
      } else if(f.type == FilterType.NEED_COLUMNS){
        m["type"] = "NEED_COLUMNS";
        m["value"] = f.value;

      } else if(f.type == FilterType.JOIN){
        m["type"] = "JOIN";
        JoinCondiction c = f.value;
        Map<String,dynamic> joinCondiction = new Map<String,dynamic>();
        joinCondiction["tableName"] = c.tableName;
        if(c.type == JoinType.INNER){
          joinCondiction["type"] = "INNER";
        } else if(c.type == JoinType.LEFT){
          joinCondiction["type"] = "LEFT";
        }

        if(c.where != null) {
          joinCondiction["where"] = c.where;
        }

        if(c.whereBindingValues != null){
          joinCondiction["whereBindingValues"] = c.whereBindingValues;
        }

        if(c.needColumns != null){
          joinCondiction["needColumns"] = c.needColumns;
        }

        if(c.matchColumns != null){
          joinCondiction["matchColumns"] = c.matchColumns;
        }
        m["value"] = joinCondiction;
      }
      params.add(m);
    });
    return params;
  }
}





class Field {

  final FieldType type;

  bool unique;

  int maxLength;

  bool primaryKey;

  bool foreignKey;

  bool autoIncrement;

  String to;

  bool index;

  Field(this.type,{this.autoIncrement = false, this.unique = false, this.maxLength = 0, this.foreignKey = false, this.primaryKey = false, this.index = false, this.to = ""});

  List<dynamic> getParamsMap(){
    List<dynamic> l = new List<dynamic>();
    if(type == FieldType.Integer){
      l.add("IntegerField");
    } else if(type == FieldType.Real){
      l.add("RealField");
    } else if(type == FieldType.Blob){
      l.add("BlobField");
    } else if(type == FieldType.Char){
      l.add("CharField");
    } else if(type == FieldType.Text){
      l.add("TextField");
    } else if(type == FieldType.Boolean){
      l.add("BooleandField");
    }

    Map<String,dynamic> params = new Map<String,dynamic>();

    if(unique) {
      params["unique"] = true;
    }


    if(maxLength != null && maxLength > 0){
      params["max_length"] = maxLength;
    }

    if(primaryKey){
      params["primary_key"] = maxLength;
    }

    if(index){
      params["index"] = true;
    }

    if(foreignKey){
      params["foreign_key"] = true;
    }

    if(autoIncrement){
      params["auto_increment"] = true;
    }

    if(to != null && to.length>0){
      params["to"] = to;
    }

    l.add(params);

    return l;
  }

}