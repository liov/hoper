package api

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

//有swagger,有没有必要做
func doc(modName string) string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "# %s接口文档  \n", modName)
	fmt.Fprint(buf, "----------  \n")
	for k, v := range svcs {
		fmt.Fprintf(buf, "# %s  \n", v.Describe())
		fmt.Fprint(buf, "----------  \n")
		value := reflect.ValueOf(v)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Type.NumIn() < 2 || method.Type.NumOut() != 2 {
				continue
			}
			methodInfo := getMethodInfo(value.Method(j))
			mName, version := parseMethodName(method.Name)
			methodInfo.path = "/api/v" + strconv.Itoa(version) + "/" + k + "/" + mName
			//title
			if methodInfo.deprecated != nil {
				fmt.Fprintf(buf, "## ~~%s-v%d(废弃)(`%s`)~~  \n", methodInfo.title, version, methodInfo.path)
			} else {
				fmt.Fprintf(buf, "## %s-v%d(`%s`)  \n", methodInfo.title, version, methodInfo.path)
			}
			//api
			fmt.Fprintf(buf, "**%s** `%s` _(Principal %s)_  \n", methodInfo.method, methodInfo.path, methodInfo.getPrincipal())

			fmt.Fprint(buf, "### 接口记录  \n")
			fmt.Fprint(buf, "|版本|操作|时间|负责人|日志|  \n")
			fmt.Fprint(buf, "| :----: | :----: | :----: | :----: | :----: |  \n")
			fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", methodInfo.createlog.version, "创建", methodInfo.createlog.date, methodInfo.createlog.auth, methodInfo.createlog.log)
			if len(methodInfo.changelog) != 0 || methodInfo.deprecated != nil {
				for _, clog := range methodInfo.changelog {
					fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", clog.version, "变更", clog.date, clog.auth, clog.log)
				}
				if methodInfo.deprecated != nil {
					fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", methodInfo.deprecated.version, "删除", methodInfo.deprecated.date, methodInfo.deprecated.auth, methodInfo.deprecated.log)
				}
			}

			fmt.Fprint(buf, "### 参数信息  \n")
			if method.Type.NumIn() == 3 {
				fmt.Fprint(buf, "|字段名称|字段类型|字段描述|校验要求|  \n")
				fmt.Fprint(buf, "| :----: | :----: | :----: | :----: |  \n")
				params := getParamTable(method.Type.In(2).Elem(), "")
				for i := range params {
					fmt.Fprintf(buf, "|%s|%s|%s|%s|  \n", params[i].json, params[i].typ, params[i].annotation, params[i].validator)
				}

			} else {
				fmt.Fprint(buf, "无需参数")
			}
		}
	}
	return buf.String()
}

type ParamTable struct {
	json, annotation, typ, validator string
}

func getParamTable(param reflect.Type, pre string) []*ParamTable {
	param = reflect3.OriginalType(param)

	var res []*ParamTable
	for i := 0; i < param.NumField(); i++ {
		/*		if param.AssignableTo(reflect.TypeOf(response.File{})) {
				return "下载文件"
			}*/
		var p ParamTable
		field := param.Field(i)
		p.json = strings.Split(field.Tag.Get("json"), ",")[0]
		if p.json == "-" {
			continue
		}
		if p.json == "" {
			p.json = pre + strings2.ConvertToCamelCase(p.json)
		} else {
			p.json = pre + p.json
		}
		p.annotation = field.Tag.Get("annotation")
		if p.annotation == "-" {
			p.annotation = p.json
		}
		p.typ = getJsType(field.Type)
		res = append(res, &p)
		if p.typ == "object" || p.typ == "[]object" {
			sub := getParamTable(field.Type, p.json+".")
			res = append(res, sub...)
		}
	}
	return res
}

func getJsType(typ reflect.Type) string {
	t := time.Time{}
	if typ == reflect.TypeOf(t) || typ == reflect.TypeOf(&t) {
		return "date"
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Array, reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			return "string"
		}
		return "[]" + getJsType(typ.Elem())
	case reflect.Ptr:
		return getJsType(typ.Elem())
	case reflect.Struct:
		return "object"
	}
	return "string"
}
