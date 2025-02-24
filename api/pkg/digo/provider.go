package digo

import (
	aflog "github.com/habiliai/alice/api/pkg/log"
	"github.com/pkg/errors"
)

type (
	ObjectKey    string
	RegisterFunc = func(container *Container) (any, error)
)

var (
	objectProviders = map[ObjectKey]RegisterFunc{}
)

func Register(key ObjectKey, provider RegisterFunc) {
	objectProviders[key] = provider
}

func Get[T any](c *Container, key ObjectKey) (T, error) {
	var (
		svc T
	)
	if untypedSvc, ok := c.objects[key]; ok {
		svc, ok = untypedSvc.(T)
		if !ok {
			return svc, errors.Errorf("type miss matched. typeof untypedSvc: %T, typeof svc: %T", untypedSvc, svc)
		}
		return svc, nil
	}

	provider, ok := objectProviders[key]
	if !ok {
		return svc, errors.New("service provider not found")
	}

	untypedSvc, err := provider(c)
	if err != nil {
		return svc, err
	}
	logger.Info("succeeded to create service.", "key", key, "svc", untypedSvc)

	c.objects[key] = untypedSvc

	svc, ok = untypedSvc.(T)
	if !ok {
		return svc, errors.Errorf("type miss matched. typeof untypedSvc: %T, typeof svc: %T", untypedSvc, svc)
	}

	return svc, nil
}

func MustGet[T any](c *Container, key ObjectKey) T {
	svc, err := Get[T](c, key)
	if err != nil {
		logger.Error("failed to get service.", "key", key, aflog.Err(err))
		panic(err)
	}

	return svc
}

func Set(c *Container, key ObjectKey, svc any) {
	c.objects[key] = svc
}
