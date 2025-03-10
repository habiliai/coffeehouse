package di

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

func Get[T any](ctx context.Context, key ObjectKey) (res T, err error) {
	c, ok := ctx.Value(containerCtxKey).(*Container)
	if !ok {
		err = errors.New("container not found")
		return
	}

	res, ok = c.objects[key].(T)
	if ok {
		return
	}

	fn, ok := registry[key]
	if !ok {
		err = errors.Errorf("object %s not registered", key)
		return
	}

	obj, err := fn(ctx, c.Env)
	if err != nil {
		return
	}

	c.objects[key] = obj

	res, ok = obj.(T)
	if !ok {
		err = errors.Errorf("object %s is not of type %T", key, res)
		return
	}

	return
}

func MustGet[T any](c context.Context, key ObjectKey) T {
	res, err := Get[T](c, key)
	if err != nil {
		panic(fmt.Sprintf("error: %+v", err))
	}

	return res
}

func Set[T any](ctx context.Context, key ObjectKey, obj T) {
	c, ok := ctx.Value(containerCtxKey).(*Container)
	if !ok {
		panic("container not found")
	}

	c.objects[key] = obj
}
