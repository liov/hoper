package constraints

import "context"

type GRPCServiceMethod[REQ, RES any] func(context.Context, REQ) (RES, error)
