package cache

import (
	"context"
	"github.com/eine-doodka/twoStage/customerrors"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

type Cache interface {
	Set(ctx context.Context, uuid uuid.UUID, code string) error
	Get(ctx context.Context, uuid uuid.UUID) (string, error)
}

type Impl struct {
	redis  *redis.Client
	expire time.Duration
}

func NewImpl(redis *redis.Client, expire time.Duration) *Impl {
	return &Impl{
		redis:  redis,
		expire: expire,
	}
}

func (i *Impl) Set(ctx context.Context, uuid uuid.UUID, code string) error {
	return i.redis.Set(ctx, uuid.String(), code, i.expire).Err()
}

func (i *Impl) Get(ctx context.Context, uuid uuid.UUID) (string, error) {
	result, err := i.redis.Get(ctx, uuid.String()).Result()
	if err == redis.Nil {
		return "", customerrors.ErrNotFound
	} else if err != nil {
		return "", err
	}
	return result, nil
}
