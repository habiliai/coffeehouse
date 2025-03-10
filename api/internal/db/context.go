package db

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	sessionCtxKey = uuid.New()
)

func OpenSession(ctx context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
	tx, ok := ctx.Value(sessionCtxKey).(*gorm.DB)
	if ok {
		return ctx, tx
	}

	return WithSession(ctx, db)
}

func WithSession(ctx context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
	tx := db.WithContext(ctx)

	return context.WithValue(ctx, sessionCtxKey, tx), tx
}
