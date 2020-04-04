package gql

import (
	"io/ioutil"
	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/utils/log"
)

var FilePath = "../protobuf/gql/"

//简直浪费时间，resolvable.makeScalarExec把类型定义固定，为了所谓graphql规范，简直垃圾，
//js最大整数值明明是Math.pow(2, 53) - 1     // 9007199254740991
//凭什么限定i32，况且js还有bigNumber
//浪费时间，知道这库为啥star不多了，垃圾
func Graphql(app *iris.Application, filePath, modName string, resolver interface{}) {
	FilePath = filePath
	f, err := os.Open(FilePath + modName + "/" + modName + ".service.pb.graphqls")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	schema := graphql.MustParseSchema(string(data), resolver,
		graphql.UseStringDescriptions(), graphql.UseFieldResolvers())
	app.Post("/api/graphql", iris.FromStd(&relay.Handler{Schema: schema}))
}

/*func GraphqlV2(app *iris.Application, filePath, modName string, resolver interface{}) {
	FilePath = filePath
	f, err := os.Open(FilePath + modName + "/" + modName + ".service.pb.graphqls")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	var sources = []*ast.Source{
		{Name: "schema.graphql", Input: string(data), BuiltIn: false}}
	var parsedSchema = gqlparser.MustLoadSchema(sources...)

	srv := handler.NewDefaultServer(svc)
	app.Post("/api/graphql", iris.FromStd(srv))
}
*/
