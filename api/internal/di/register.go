package di

import "context"

type (
	RegisterFn func(ctx context.Context, env Env) (any, error)
)

var (
	registry = map[ObjectKey]RegisterFn{}
)

func Register(key ObjectKey, fn RegisterFn) {
	registry[key] = fn
}
