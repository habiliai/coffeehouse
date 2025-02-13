package helpers

import (
	"context"
	"gorm.io/gorm"
)

type contextKey string

const (
	contextKeyGithubToken contextKey = "_githubToken"
	contextKeyAuthToken   contextKey = "_authToken"
	contextKeyTx          contextKey = "_tx"
	contextKeyDeviceId    contextKey = "_deviceId"
)

func GetGithubToken(ctx context.Context) string {
	value, ok := ctx.Value(contextKeyGithubToken).(string)
	if !ok {
		return ""
	}

	return value
}

func WithGithubToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, contextKeyGithubToken, token)
}

func GetAuthToken(ctx context.Context) string {
	value, ok := ctx.Value(contextKeyAuthToken).(string)
	if !ok {
		return ""
	}

	return value
}

func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, contextKeyAuthToken, token)
}

func GetTx(ctx context.Context) *gorm.DB {
	value, ok := ctx.Value(contextKeyTx).(*gorm.DB)
	if !ok {
		return nil
	}

	return value
}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, contextKeyTx, tx)
}

func GetDeviceId(ctx context.Context) string {
	value, ok := ctx.Value(contextKeyDeviceId).(string)
	if !ok {
		return ""
	}

	return value
}

func WithDeviceId(ctx context.Context, deviceId string) context.Context {
	return context.WithValue(ctx, contextKeyDeviceId, deviceId)
}
