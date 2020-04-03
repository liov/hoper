package gql

import (
	"context"
	"strconv"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/kataras/iris/v12"
	model "github.com/liov/hoper/go/v2/protobuf/user"
)

func GraphqlRouter(app *iris.Application) {

	s := `
                schema {
                      query: Query
                }
				type Query {
  					getUser(id: ID!): User
				}
                type User {
                        Id:ID!
						Name:String!
						Gender:Gender!
						Phone:String!
                }
				enum Gender{
					男
					女
					未填
				}
        `

	schema := graphql.MustParseSchema(s, &Resolver{}, graphql.UseStringDescriptions())

	graphqlRouter := app.Party("/api/v1/graphql")
	{
		graphqlRouter.Post("", iris.FromStd(&relay.Handler{Schema: schema}))
	}
}

type UserResolver struct {
	user model.User
}

func (u *UserResolver) Id() graphql.ID       { return graphql.ID(strconv.Itoa(int(u.user.Id))) }
func (u *UserResolver) Name() string         { return u.user.Name }
func (u *UserResolver) Gender() model.Gender { return u.user.Gender }
func (u *UserResolver) Phone() string        { return u.user.Phone }

type Resolver struct {
}

func (r *Resolver) GetUser(ctx context.Context, args struct{ ID graphql.ID }) (*UserResolver, error) {
	s := UserResolver{
		user: model.User{Id: 1},
	}
	return &s, nil
}
