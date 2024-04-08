package storage

import (
	"context"
	"time"
)

type CacheInterface interface {
	Code() CodeCacheInterface
}

type CodeCacheInterface interface {
	SetCode(ctx context.Context, email, code string, expire time.Time) error
	GetCode(ctx context.Context, email string) (string, time.Time, error)
	DeleteCode(ctx context.Context, key string) error
}
