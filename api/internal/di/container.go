package di

import (
	"context"
	"github.com/google/uuid"
)

type (
	Env       string
	Container struct {
		objects map[ObjectKey]any
		Env     Env
	}
	ObjectKey uuid.UUID
)

const (
	EnvProd Env = "prod"
	EnvTest Env = "test"
)

var (
	containerCtxKey = uuid.New()
)

func NewKey() ObjectKey {
	return ObjectKey(uuid.New())
}

func WithContainer(ctx context.Context, env Env) context.Context {
	container := &Container{
		Env:     env,
		objects: map[ObjectKey]any{},
	}

	return context.WithValue(ctx, containerCtxKey, container)
}
