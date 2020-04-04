package graphql

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
)

func GraphqlRouterV2(app *iris.Application) {
	// Schema
	user := graphql.NewObject(graphql.ObjectConfig{Name: "User", Fields: graphql.Fields{
		"Name":   &graphql.Field{Type: graphql.String},
		"Id":     &graphql.Field{Type: graphql.String},
		"Gender": &graphql.Field{Type: graphql.String},
		"Phone":  &graphql.Field{Type: graphql.String},
	}})
	query := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{
		"getUser": &graphql.Field{
			Type: user,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.ID,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				/*idStr, _ := p.Args["id"].(string)
				id,_ :=strconv.ParseUint(idStr,10,64)*/
				return &model.User{Id: 1, Name: "test"}, nil
			},
		},
	}})
	schemaConfig := graphql.SchemaConfig{Query: query}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	app.Post("/api/v2/graphql", func(context context.Context) {

		var p struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(context.Request().Body).Decode(&p); err != nil {
			http.Error(context.ResponseWriter(), err.Error(), http.StatusBadRequest)
			return
		}
		params := graphql.Params{Schema: schema,
			RequestString:  p.Query,
			VariableValues: p.Variables,
			OperationName:  p.OperationName,
		}
		r := graphql.Do(params)
		context.JSON(r)
	})
}
