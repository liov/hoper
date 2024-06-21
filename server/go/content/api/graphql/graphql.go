package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/liov/hoper/server/go/content/service"
	"github.com/liov/hoper/server/go/protobuf/content"
)

type Resolver struct{}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() content.MutationResolver {
	return &mutationResolver{
		ActionServiceResolvers:  content.ActionServiceResolvers{Service: &service.ActionService{}},
		ContentServiceResolvers: content.ContentServiceResolvers{Service: &service.ContentService{}},
		MomentServiceResolvers:  content.MomentServiceResolvers{Service: &service.MomentService{}},
		DiaryServiceResolvers:   content.DiaryServiceResolvers{Service: &service.DiaryService{}},
		NoteServiceResolvers:    content.NoteServiceResolvers{Service: &service.NoteService{}},
		TestServiceResolvers:    content.TestServiceResolvers{Service: &service.TestService{}},
	}
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() content.QueryResolver {
	return &queryResolver{
		ActionServiceResolvers:  content.ActionServiceResolvers{Service: &service.ActionService{}},
		ContentServiceResolvers: content.ContentServiceResolvers{Service: &service.ContentService{}},
		MomentServiceResolvers:  content.MomentServiceResolvers{Service: &service.MomentService{}},
		DiaryServiceResolvers:   content.DiaryServiceResolvers{Service: &service.DiaryService{}},
		NoteServiceResolvers:    content.NoteServiceResolvers{Service: &service.NoteService{}},
		TestServiceResolvers:    content.TestServiceResolvers{Service: &service.TestService{}},
	}
}

type mutationResolver struct {
	content.ActionServiceResolvers
	content.ContentServiceResolvers
	content.MomentServiceResolvers
	content.DiaryServiceResolvers
	content.NoteServiceResolvers
	content.TestServiceResolvers
}
type queryResolver = mutationResolver

func (q queryResolver) TestServiceGc(ctx context.Context, in *content.GCReq) (*bool, error) {
	return q.TestServiceGC(ctx, in)
}

func NewExecutableSchema() graphql.ExecutableSchema {
	directive := func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		return next(ctx)
	}
	return content.NewExecutableSchema(content.Config{
		Resolvers: &Resolver{},
		Directives: content.DirectiveRoot{
			ActionService:  directive,
			ContentService: directive,
			DiaryService:   directive,
			MomentService:  directive,
			NoteService:    directive,
			TestService:    directive,
		},
		Complexity: content.ComplexityRoot{},
	})
}
