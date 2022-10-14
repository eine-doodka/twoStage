package cache

import (
	"context"
	"github.com/eine-doodka/twoStage/customerrors"
	"github.com/google/uuid"
)

type InMemoryCache struct {
	cache map[uuid.UUID]string
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[uuid.UUID]string),
	}
}

func (i *InMemoryCache) Set(_ context.Context, uuid uuid.UUID, code string) error {
	i.cache[uuid] = code
	return nil
}

func (i *InMemoryCache) Get(_ context.Context, uuid uuid.UUID) (string, error) {
	res, ok := i.cache[uuid]
	if !ok {
		return "", customerrors.ErrNotFound
	}
	return res, nil
}
